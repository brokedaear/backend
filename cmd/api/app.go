// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"log/slog"
)

type app struct {
	config   *config
	logger   *slog.Logger
	services *service.Service
}

// config defines application configuration.
type config struct {
	// port access for host.
	port int

	// Runtime environment, either "development", "staging", or "production".
	env Environment

	// version number, based on environment variable.
	version string

	// db is the database configuration config.
	db dbConfig

	// limiter is for limiter information for rate limiting.
	limiter struct {
		// rps is requests per second.
		rps float64

		// burst is how many bursts are allowed.
		burst int

		// enabled either disables or enables rate limited altogether.
		enabled bool
	}

	// Store secrets as bytes for easy overwrite capability.
	secrets struct {
		stripePrivateKey    []byte
		stripePublicKey     []byte
		stripeWebhookSecret []byte
	}
}

type dbConfig struct {
	// Driver is the SQL driver, like postgreSQL..
	Driver string

	// DSN is Data Source Name.
	Dsn string

	// Database parameters. On the app layer, these are found in env variables.

	Name     string
	Username string
	Password string
	Host     string
	Port     string
	SslMode  string

	// Connection rate limiting logic.

	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Environment uint8

const (
	ENV_DEVELOPMENT Environment = 0
	ENV_STAGING     Environment = 1
	ENV_PRODUCTION  Environment = 2
)
