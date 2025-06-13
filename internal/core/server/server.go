// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
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

func NewBase(logger Logger, config *Config) (*Base, error) {
	return &Base{
		logger: logger,
		config: config,
	}, nil
}
