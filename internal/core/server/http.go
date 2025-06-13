// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HttpServer represents an HTTP server that is capable of accepting routes
// and connections from clients via HTTP. Implementations for this interface
// include endpoints that serve GraphQL or simple barebones requests.
//
// This interface also implements the io.Closer interface, for use in global
// teardown operations.
type HttpServer interface {
	ListenAndServe(context.Context) error
	RegisterRoutes(...HttpRoute)
	io.Closer
}

type httpServer struct {
	*Base
	srv *http.Server
}

// NewHttpServer creates a new HTTP server using a logger and a config.
// The server comes with telemetry enabled by default.
func NewHttpServer(logger Logger, config *Config) (HttpServer, error) {
	s := &httpServer{}

	b, err := NewBase(logger, config)
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(config.Addr.String(), config.Port.String())

	b.listener, err = net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s.Base = b

	s.srv = &http.Server{
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s, nil
}

// ListenAndServe listens to specified route endpoints given by route functions
// specified by the function signature. The server is terminated via error or
// the server's interface io.Closer method. An error is only returned when the
// closure results from an error.
func (s httpServer) ListenAndServe(ctx context.Context) error {
	var serverError error

	serverCtx, serverCancel := context.WithCancel(ctx)

	defer serverCancel()

	// routes = append(routes, routeHealthcheck)
	go func() {
		defer serverCancel()
		err := s.srv.Serve(s.listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(err.Error())
			serverError = err
		}
	}()

	<-serverCtx.Done()

	return serverError
}

func (s httpServer) Close() error {
	shutdownCtx, shutdownCancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))

	defer shutdownCancel()

	err := s.srv.Shutdown(shutdownCtx)
	if err != nil {
		s.logger.Warn("failed to shutdown http server, killing", "err", err)
		err = s.srv.Close()
		if err != nil {
			return err
		}
	}

	err = s.listener.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close http listener")
	}

	s.logger.Info("http server closed")

	return nil
}

type HttpRoute interface {
	String() string
	Route() http.HandlerFunc
}

func (s httpServer) RegisterRoutes(routes ...HttpRoute) {
	s.srv.Handler = s.registerRoutes(routes...)
}

// registerRoutes needs to be REFACTORED. TODO.
func (s httpServer) registerRoutes(routes ...HttpRoute) http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc http.HandlerFunc) {
		mux.HandleFunc(pattern, handlerFunc)
	}

	if s.config.Telemetry {
		handleFunc = func(pattern string, handlerFunc http.HandlerFunc) {
			h := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
			mux.Handle(pattern, h)
		}
	}

	for _, v := range routes {
		handleFunc(v.String(), v.Route())
	}

	if s.config.Telemetry {
		return otelhttp.NewHandler(mux, "/")
	}

	return mux
}
