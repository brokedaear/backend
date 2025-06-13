// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"backend.brokedaear.com/internal/core/server"
)

type monitorServer struct {
	server.HttpServer
}

func newMonitorServer(logger server.Logger, config *server.Config) (*monitorServer, error) {
	s, err := server.NewHttpServer(logger, config)
	if err != nil {
		return nil, err
	}

	return &monitorServer{
		HttpServer: s,
	}, nil
}

func (s *monitorServer) Start(ctx context.Context) error {
	return s.ListenAndServe(ctx)
}
