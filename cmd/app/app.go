// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"backend.brokedaear.com/internal/core/server"
)

type appServer struct {
	server.HTTPServer
	logger server.Logger
}

func newAppServer(
	ctx context.Context,
	logger server.Logger,
	config *server.Config,
) (*appServer, error) {
	s, err := server.NewHTTPServer(ctx, logger, config)
	if err != nil {
		return nil, err
	}

	return &appServer{
		HTTPServer: s,
		logger:     logger,
	}, nil
}

func (s *appServer) Start(ctx context.Context) error {
	return s.ListenAndServe(ctx)
}

func (s *appServer) Close() error {
	return s.HTTPServer.Close()
}
