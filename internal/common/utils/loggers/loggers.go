// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// logger implements builder functions for loggers that implement the logging
// interface in common/utils.

package loggers

import (
	"log/slog"

	"go.uber.org/zap"
)

// ZapLogger is a facade for zap.Logger to implement the Logger interface.
type ZapLogger struct {
	zap     *zap.Logger
	sugared *zap.SugaredLogger
}

// NewZap creates a new instance of a Development Zap logger.
func NewZap() (*ZapLogger, error) {
	z, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		zap:     z,
		sugared: z.Sugar(),
	}, nil
}

func (l *ZapLogger) Info(msg string, args ...any) {
	l.sugared.Infof(msg, args...)
}

func (l *ZapLogger) Debug(msg string, args ...any) {
	l.sugared.Debugf(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...any) {
	l.sugared.Warnf(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...any) {
	l.sugared.Errorf(msg, args...)
}

func (l *ZapLogger) Sync() error {
	return l.zap.Sync()
}

// NewZapProd creates a zap logger. It also returns the logger's flushing
// method.
func NewZapProd() (*zap.Logger, func() error) {
	l, _ := zap.NewProduction()
	return l, l.Sync
}

// NewPrettySlog creates a logger using the stdlib `slog` package.
func NewPrettySlog() *slog.Logger {
	slogHandlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	return slog.New(New(slogHandlerOptions))
}
