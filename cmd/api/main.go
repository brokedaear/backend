// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"backend.brokedaear.com/internal/common/utils/loggers"
	"backend.brokedaear.com/internal/core/service"
)

const version = "1.0.0"

func main() {
	envFlag := flag.Int("environment", 0, "application environment 0-2")

	flag.Parse()

	if *envFlag < 0 || *envFlag > 2 {
		panic("invalid application environment")
	}

	config := &config{
		port:    7402,
		env:     Environment(*envFlag),
		version: version,
		secrets: struct {
			stripePrivateKey    []byte
			stripePublicKey     []byte
			stripeWebhookSecret []byte
		}{},
	}

	slogHandlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if config.env == ENV_DEVELOPMENT || config.env == ENV_STAGING {
		slogHandlerOptions.AddSource = true
		slogHandlerOptions.Level = slog.LevelDebug
	}

	logger := slog.New(loggers.NewHandler(slogHandlerOptions))

	logger.Info(fmt.Sprintf("application environment set to %d", config.env))

	services := service.NewServices()

	app := &app{
		config:   config,
		logger:   logger,
		services: services,
	}

	err := app.server()

	logger.Error(err.Error())
	os.Exit(1)
}
