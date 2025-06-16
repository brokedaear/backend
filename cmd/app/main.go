// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"backend.brokedaear.com"
	"backend.brokedaear.com/internal/common/infra"
	"backend.brokedaear.com/internal/common/utils/loggers"
	"backend.brokedaear.com/internal/core/server"
)

const (
	version     = "0.1.0"
	environment = "development"
	port        = 1025
	address     = "localhost"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("application error: %v\n", err) //nolint:forbidigo // structured logger not available on init failure
		os.Exit(1)
	}
}

func run() error {
	env, err := backend.EnvFromString(environment)
	if err != nil {
		return fmt.Errorf("failed to parse environment: %w", err)
	}

	config := &loggers.ZapConfig{
		Env:                env,
		OtelServiceName:    "app",
		OtelLoggerProvider: nil,
		CustomZapper:       nil,
		WithTelemetry:      false,
	}
	logger, err := loggers.NewZap(config)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			fmt.Printf("failed to sync zap logger: %v\n", syncErr) //nolint:forbidigo // logger sync failure, no alternative
		}
	}()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	defer cancel()

	cfg, err := server.NewConfig(address, port, environment, version)
	if err != nil {
		logger.Error("failed to initialize config", "error", err)
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	s, err := newAppServer(ctx, logger, cfg)
	if err != nil {
		logger.Error("failed to create monitor server", "error", err)
		return fmt.Errorf("failed to create monitor server: %w", err)
	}

	return runService(ctx, logger, s)
}

func runService(ctx context.Context, logger server.Logger, s *appServer) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(
		func() error {
			logger.Info("starting monitor server", "address", address, "port", port)
			return s.Start(gCtx)
		},
	)

	g.Go(
		func() error {
			<-gCtx.Done()
			logger.Info("shutdown signal received, initiating graceful shutdown")

			const shutdownTimeout = 30 * time.Second
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer shutdownCancel()

			select {
			case <-shutdownCtx.Done():
				logger.Warn("shutdown timeout exceeded, forcing exit")
			default:
			}

			if teardownErr := infra.Teardown(s); teardownErr != nil {
				logger.Error("error during teardown", "error", teardownErr)
				return teardownErr
			}

			logger.Info("graceful shutdown completed")
			return nil
		},
	)

	if err := g.Wait(); err != nil {
		logger.Error("monitor service error", "error", err)
		return fmt.Errorf("monitor service error: %w", err)
	}

	return nil
}
