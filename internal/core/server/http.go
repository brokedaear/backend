// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

// HttpServer represents an HTTP server that is capable of accepting routes
// and connections from clients via HTTP. Implementations for this interface
// include endpoints that serve GraphQL or simple barebones requests.
//
// This interface also implements the io.Closer interface, for use in global
// teardown operations.
type HttpServer interface {
	ListenAndServe(context.Context, ...func(*http.ServeMux)) error
	io.Closer
}

type httpServer struct {
	*Base
}

// NewHttpServer creates a new HTTP server using a logger and a config.
func NewHttpServer(logger Logger, config *Config) (HttpServer, error) {
	b, err := NewBase(logger, config)
	if err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", config.Addr.String())
	if err != nil {
		return nil, err
	}

	b.listener = lis

	return &httpServer{
		Base: b,
	}, nil
}

// ListenAndServe listens to specified route endpoints given by route functions
// specified by the function signature. The server is terminated via error or
// the server's interface io.Closer method. An error is only returned when the
// closure results from an error.
func (s httpServer) ListenAndServe(ctx context.Context, routes ...func(*http.ServeMux)) error {
	var serverError error

	serverCtx, serverCancel := context.WithCancel(ctx)

	defer serverCancel()

	handler := http.NewServeMux()

	// routes = append(routes, routeHealthcheck)

	for _, addRoute := range routes {
		addRoute(handler)
	}

	go func() {
		defer serverCancel()
		err := http.Serve(s.listener, handler)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(err.Error())
			serverError = err
		}
	}()

	<-serverCtx.Done()

	return serverError
}

// Close closes the server.
func (s httpServer) Close() error {
	err := s.listener.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close http server")
	}

	s.logger.Info("http server closed")

	return nil
}
