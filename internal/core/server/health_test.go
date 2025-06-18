// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"backend.brokedaear.com/internal/core/server"
	"backend.brokedaear.com/pkg/assert"
	"backend.brokedaear.com/pkg/test"
)

func TestNewHealthServer(t *testing.T) {
	logger := test.NewMockLogger()

	healthServer := server.NewHealthServer(logger)
	assert.NotEqual(t, healthServer, nil)
}

func TestHealthServer_SetServingStatus(t *testing.T) {
	tests := []struct {
		test.CaseBase
		service string
		status  grpc_health_v1.HealthCheckResponse_ServingStatus
	}{
		{
			CaseBase: test.NewCaseBase(
				"set overall server status to serving",
				nil,
				false,
			),
			service: "",
			status:  grpc_health_v1.HealthCheckResponse_SERVING,
		},
		{
			CaseBase: test.NewCaseBase(
				"set specific service status to not serving",
				nil,
				false,
			),
			service: "myservice",
			status:  grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		},
		{
			CaseBase: test.NewCaseBase(
				"set service status to unknown",
				nil,
				false,
			),
			service: "unknown-service",
			status:  grpc_health_v1.HealthCheckResponse_UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				_ = t
				logger := test.NewMockLogger()
				healthServer := server.NewHealthServer(logger)

				// This should not panic.
				healthServer.SetServingStatus(tt.service, tt.status)
			},
		)
	}
}

func TestHealthServer_Check(t *testing.T) {
	tests := []struct {
		test.CaseBase
		service      string
		setupService string
		setupStatus  grpc_health_v1.HealthCheckResponse_ServingStatus
		expectedCode codes.Code
	}{
		{
			CaseBase: test.NewCaseBase(
				"check existing service returns status",
				grpc_health_v1.HealthCheckResponse_SERVING,
				false,
			),
			service:      "myservice",
			setupService: "myservice",
			setupStatus:  grpc_health_v1.HealthCheckResponse_SERVING,
			expectedCode: codes.OK,
		},
		{
			CaseBase: test.NewCaseBase(
				"check overall server status",
				grpc_health_v1.HealthCheckResponse_NOT_SERVING,
				false,
			),
			service:      "",
			setupService: "",
			setupStatus:  grpc_health_v1.HealthCheckResponse_NOT_SERVING,
			expectedCode: codes.OK,
		},
		{
			CaseBase: test.NewCaseBase(
				"check unknown service returns not found",
				nil,
				true,
			),
			service:      "unknown",
			setupService: "known",
			setupStatus:  grpc_health_v1.HealthCheckResponse_SERVING,
			expectedCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				ctx := t.Context()
				logger := test.NewMockLogger()
				healthServer := server.NewHealthServer(logger)

				// Setup the service status.
				healthServer.SetServingStatus(tt.setupService, tt.setupStatus)

				req := &grpc_health_v1.HealthCheckRequest{
					Service: tt.service,
				}

				resp, err := healthServer.Check(ctx, req)
				assert.ErrorOrNoError(t, err, tt.WantErr)

				if tt.WantErr {
					st, ok := status.FromError(err)
					assert.Equal(t, ok, true)
					assert.Equal(t, st.Code(), tt.expectedCode)
				} else {
					assert.NotEqual(t, resp, nil)
					assert.Equal(
						t,
						resp.GetStatus(),
						tt.Want.(grpc_health_v1.HealthCheckResponse_ServingStatus),
					)
				}
			},
		)
	}
}

func TestHealthServer_Watch(t *testing.T) {
	ctx := t.Context()
	logger := test.NewMockLogger()
	healthServer := server.NewHealthServer(logger)

	service := "testservice"
	initialStatus := grpc_health_v1.HealthCheckResponse_SERVING

	// Setup the service.
	healthServer.SetServingStatus(service, initialStatus)

	// Create a mock stream.
	mockStream := &mockHealthWatchServer{
		ctx:      ctx,
		cancel:   nil,
		received: make(chan *grpc_health_v1.HealthCheckResponse, 10),
	}

	req := &grpc_health_v1.HealthCheckRequest{
		Service: service,
	}

	// Start watching in a goroutine.
	watchDone := make(chan error, 1)
	go func() {
		watchDone <- healthServer.Watch(req, mockStream)
	}()

	// Should receive initial status immediately.
	select {
	case resp := <-mockStream.received:
		assert.Equal(t, resp.GetStatus(), initialStatus)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("did not receive initial health status")
	}

	// Change status and should receive update.
	newStatus := grpc_health_v1.HealthCheckResponse_NOT_SERVING
	healthServer.SetServingStatus(service, newStatus)

	select {
	case resp := <-mockStream.received:
		assert.Equal(t, resp.GetStatus(), newStatus)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("did not receive health status update")
	}

	// Cancel context to end watch.
	mockStream.cancel()

	select {
	case err := <-watchDone:
		assert.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("watch did not complete after context cancellation")
	}
}

func TestHealthServer_WatchUnknownService(t *testing.T) {
	ctx := t.Context()
	logger := test.NewMockLogger()
	healthServer := server.NewHealthServer(logger)

	mockStream := &mockHealthWatchServer{
		ctx:      ctx,
		cancel:   nil,
		received: make(chan *grpc_health_v1.HealthCheckResponse, 10),
	}

	req := &grpc_health_v1.HealthCheckRequest{
		Service: "unknown",
	}

	err := healthServer.Watch(req, mockStream)
	assert.ErrorOrNoError(t, err, true)

	st, ok := status.FromError(err)
	assert.Equal(t, ok, true)
	assert.Equal(t, st.Code(), codes.NotFound)
}

func TestHealthServer_Shutdown(t *testing.T) {
	logger := test.NewMockLogger()
	healthServer := server.NewHealthServer(logger)

	// Setup some services.
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("service1", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("service2", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	// Shutdown should not panic.
	healthServer.Shutdown()

	// After shutdown, all services should report NOT_SERVING.
	ctx := t.Context()

	resp, err := healthServer.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
	assert.NoError(t, err)
	assert.Equal(t, resp.GetStatus(), grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	resp, err = healthServer.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "service1"})
	assert.NoError(t, err)
	assert.Equal(t, resp.GetStatus(), grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

// mockHealthWatchServer implements grpc_health_v1.Health_WatchServer for testing.
type mockHealthWatchServer struct {
	ctx      context.Context
	cancel   context.CancelFunc
	received chan *grpc_health_v1.HealthCheckResponse
}

func (m *mockHealthWatchServer) Send(resp *grpc_health_v1.HealthCheckResponse) error {
	select {
	case m.received <- resp:
		return nil
	case <-m.ctx.Done():
		return m.ctx.Err()
	}
}

func (m *mockHealthWatchServer) Context() context.Context {
	if m.cancel == nil {
		m.ctx, m.cancel = context.WithCancel(m.ctx)
	}
	return m.ctx
}

// Implement the remaining methods required by grpc.ServerStreamingServer.
func (m *mockHealthWatchServer) SendMsg(_ any) error {
	return nil
}

func (m *mockHealthWatchServer) RecvMsg(_ any) error {
	return nil
}

func (m *mockHealthWatchServer) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockHealthWatchServer) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockHealthWatchServer) SetTrailer(metadata.MD) {
}
