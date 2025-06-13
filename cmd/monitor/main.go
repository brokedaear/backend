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

	"backend.brokedaear.com/internal/common/infra"
	"backend.brokedaear.com/internal/common/utils/loggers"
	"backend.brokedaear.com/internal/core/server"
)

const (
	version     = "0.1.0"
	environment = 0
	port        = 1025
	address     = "localhost"
)

func main() {
	logger, err := loggers.NewZap()
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	defer func(logger *loggers.ZapLogger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("failed to sync zap logger: %v\n", err)
		}
	}(logger)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	defer stop()

	cfg, err := server.NewConfig(address, port, environment, version)
	if err != nil {
		logger.Warn("failed to initialize config", "msg", err)
		return
	}

	s, err := newMonitorServer(logger, cfg)
	if err != nil {
		logger.Warn("failed to create monitor server", "msg", err)
		return
	}

	var (
		errs = make(chan error)
		halt = make(chan os.Signal, 1)
	)

	signal.Notify(
		halt,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	go func() {
		defer stop()
		for {
			select {
			case err := <-errs:
				logger.Error(err.Error())
				return
			case <-halt:
				logger.Warn("Received halt, shutting down...")
				return
			}
		}
	}()

	<-ctx.Done()

	err = infra.Teardown(s)
	if err != nil {
		logger.Error(err.Error())
	}
}
