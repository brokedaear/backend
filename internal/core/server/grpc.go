// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// GRPCServer represents a gRPC server that is capable of accepting service
// registrations and connections from clients via gRPC. Implementations for
// this interface include servers that serve gRPC services with health checks
// and reflection enabled.
//
// This interface also implements the io.Closer interface, for use in global
// teardown operations.
type GRPCServer interface {
	ListenAndServe(context.Context) error
	RegisterService(desc *grpc.ServiceDesc, impl any)
	SetHealthStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus)
	io.Closer
}

type grpcServer struct {
	*Base
	srv          *grpc.Server
	healthServer *HealthServer
}

// NewGRPCServer creates a new gRPC server using a logger and a config.
// The server comes with telemetry, health checks, and reflection enabled by default.
func NewGRPCServer(ctx context.Context, logger Logger, config *Config) (GRPCServer, error) {
	b, err := NewBase(ctx, logger, config)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	address, err := config.newURIAddress()
	if err != nil {
		return nil, err
	}

	b.listener, err = net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(panicRecoveryUnaryInterceptor(b.logger)),
		grpc.StreamInterceptor(panicRecoveryStreamInterceptor(b.logger)),
	}

	if config.Telemetry {
		opts = append(
			opts,
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
	}

	srv := grpc.NewServer(opts...)

	healthServer := NewHealthServer(b.logger)
	grpc_health_v1.RegisterHealthServer(srv, healthServer)

	// Set initial overall server health to serving.
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(srv)

	return &grpcServer{
		Base:         b,
		srv:          srv,
		healthServer: healthServer,
	}, nil
}

// ListenAndServe starts the gRPC server and listens for incoming connections.
// The server is terminated via error or the server's interface io.Closer method.
// An error is only returned when the closure results from an error.
func (s grpcServer) ListenAndServe(ctx context.Context) error {
	var serverError error

	serverCtx, serverCancel := context.WithCancel(ctx)

	defer serverCancel()

	go func() {
		defer serverCancel()
		err := s.srv.Serve(s.listener)
		if err != nil {
			s.logger.Error(err.Error())
			serverError = err
		}
	}()

	<-serverCtx.Done()

	return serverError
}

// RegisterService registers a gRPC service with the server.
func (s grpcServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.srv.RegisterService(desc, impl)
}

// SetHealthStatus sets the health status for a specific service.
// Use empty string for the overall server health status.
func (s grpcServer) SetHealthStatus(service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	s.healthServer.SetServingStatus(service, status)
}

// GetListener returns the underlying listener for testing purposes.
func (s grpcServer) GetListener() net.Listener {
	return s.listener
}

func (s grpcServer) Close() error {
	const shutdownTimeout = 20 * time.Second

	// Shutdown health server first to notify watchers.
	s.healthServer.Shutdown()

	shutdownCtx, shutdownCancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(shutdownTimeout),
	)

	defer shutdownCancel()

	stopped := make(chan struct{})

	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-shutdownCtx.Done():
		s.logger.Warn("grpc server graceful shutdown timeout, forcing stop")
		s.srv.Stop()
	}

	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			s.logger.Warn("failed to close grpc listener", "err", err)
		}
	}

	s.logger.Info("grpc server closed")

	return nil
}

// panicRecoveryUnaryInterceptor returns a unary server interceptor that recovers from panics.
func panicRecoveryUnaryInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		var resp any
		var err error

		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				logger.Error(
					fmt.Sprintf("panic recovered in gRPC unary handler: %v", r),
					"stack",
					string(stack),
					"method",
					info.FullMethod,
				)
				resp = nil
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

// panicRecoveryStreamInterceptor returns a stream server interceptor that recovers from panics.
func panicRecoveryStreamInterceptor(logger Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		var err error

		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				logger.Error(
					fmt.Sprintf("panic recovered in gRPC stream handler: %v", r),
					"stack",
					string(stack),
					"method",
					info.FullMethod,
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		err = handler(srv, ss)
		return err
	}
}
