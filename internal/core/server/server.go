// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"net"

	"backend.brokedaear.com/internal/common/telemetry"
)

// Base is the base setup for any kind of server.
type Base struct {
	logger    Logger
	Telemetry telemetry.Telemetry
	config    *Config
	listener  net.Listener
}

func NewBase(ctx context.Context, logger Logger, config *Config) (*Base, error) {
	if config == nil {
		return nil, ErrNilConfig
	}

	tc := &telemetry.Config{
		ServiceName:    "server",
		ServiceVersion: config.Version.String(),
		ServiceID:      "server-1",
		ExporterConfig: *telemetry.NewExporterConfig(
			telemetry.WithType(telemetry.ExporterTypeStdout),
		),
	}

	t, err := telemetry.New(ctx, tc)
	if err != nil {
		return nil, err
	}

	return &Base{
		logger:    logger,
		config:    config,
		listener:  nil,
		Telemetry: t,
	}, nil
}

type BaseError string

func (b BaseError) Error() string {
	return string(b)
}

var ErrNilConfig BaseError = "config is nil"
