// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"backend.brokedaear.com/internal/core/server"
)

type monitorServer struct {
	server.HTTPServer
	logger server.Logger
}

func newMonitorServer(
	ctx context.Context,
	logger server.Logger,
	config *server.Config,
) (*monitorServer, error) {
	s, err := server.NewHTTPServer(ctx, logger, config)
	if err != nil {
		return nil, err
	}

	return &monitorServer{
		HTTPServer: s,
		logger:     logger,
	}, nil
}

func (s *monitorServer) Start(ctx context.Context) error {
	s.logger.Info("hello from start")
	return s.ListenAndServe(ctx)
}

func (s *monitorServer) Close() error {
	return s.HTTPServer.Close()
}
