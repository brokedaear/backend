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

	"github.com/alexliesenfeld/health"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPServer represents an HTTP server that is capable of accepting routes
// and connections from clients via HTTP. Implementations for this interface
// include endpoints that serve GraphQL or simple barebones requests.
//
// This interface also implements the io.Closer interface, for use in global
// teardown operations.
type HTTPServer interface {
	ListenAndServe(context.Context) error
	// RegisterRoutes(...HTTPRoute)
	io.Closer
}

type httpServer struct {
	*Base
	srv *http.Server
}

const httpHealthTimeout = 10 * time.Second

// NewHTTPServer creates a new HTTP server using a logger and a config.
// The server comes with telemetry enabled by default.
func NewHTTPServer(ctx context.Context, logger Logger, config *Config) (HTTPServer, error) {
	const (
		readTimeout  = 10 * time.Second
		writeTimeout = 30 * time.Second
	)

	b, err := NewBase(ctx, logger, config)
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

	mux := http.NewServeMux()

	checker := health.NewChecker(health.WithCacheDuration(1*time.Second), health.WithTimeout(httpHealthTimeout))

	mux.Handle("/health", health.NewHandler(checker))

	return &httpServer{
		Base: b,
		srv: &http.Server{
			IdleTimeout:  time.Minute,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			Handler:      mux,
		},
	}, nil
}

// ListenAndServe listens to specified route endpoints given by route functions
// specified by the function signature. The server is terminated via error or
// the server's interface io.Closer method. An error is only returned when the
// closure results from an error.
func (s httpServer) ListenAndServe(ctx context.Context) error {
	var serverError error

	serverCtx, serverCancel := context.WithCancel(ctx)

	defer serverCancel()

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
	const shutdownTimeout = 20 * time.Second
	shutdownCtx, shutdownCancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(shutdownTimeout),
	)

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

type HTTPRoute interface {
	String() string
	Route() http.HandlerFunc
}

func (s httpServer) RegisterRoutes(routes ...HTTPRoute) {
	s.srv.Handler = s.registerRoutes(routes...)
}

// registerRoutes needs to be REFACTORED.
// TODO: Refactor registerRoutes
func (s httpServer) registerRoutes(routes ...HTTPRoute) http.Handler {
	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc http.HandlerFunc) {
		mux.HandleFunc(pattern, handlerFunc)
	}

	if s.config.Telemetry {
		handleFunc = func(pattern string, handlerFunc http.HandlerFunc) {
			h := otelhttp.WithRouteTag(pattern, handlerFunc)
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
