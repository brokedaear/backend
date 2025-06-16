// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// logger implements builder functions for loggers that implement the logging
// interface in common/utils.

package loggers

import (
	"errors"
	"io"
	"net/url"
	"os"
	"syscall"

	"backend.brokedaear.com"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapConfig struct {
	Env                backend.Environment
	OtelServiceName    string
	OtelLoggerProvider *log.LoggerProvider
	CustomZapper       *CustomZapWriter
	WithTelemetry      bool
}

// ZapWriter satisfies the zap.Sink interface.
type ZapWriter interface {
	Sync() error
	io.Closer
	io.Writer
}

type CustomZapWriter struct {
	customPath string
	// customWriterKey represents the key of this writer in the repository of
	// all sinks. The key is stored in a map with the value of a sink.
	// Therefore, customWriterKey must be unique.
	customWriterKey   string
	customFunctionKey string
	zw                ZapWriter
}

func NewCustomZapWriter(customPath, customWriterKey, customFunctionKey string, zapWriter ZapWriter) *CustomZapWriter {
	return &CustomZapWriter{
		customPath:        customPath,
		customWriterKey:   customWriterKey,
		customFunctionKey: customFunctionKey,
		zw:                zapWriter,
	}
}

func (c CustomZapWriter) Validate() error {
	if c.customFunctionKey == "" || c.customPath == "" || c.customWriterKey == "" {
		return errors.New("all config fields must not be empty")
	}
	return nil
}

func (c CustomZapWriter) Value() any {
	return c
}

// NewZap creates a new instance of a Zap logger. The logger type is based on
// the runtime environment. A nil CustomZapWriter config opts out of using
// a custom zap writer.
func NewZap(config *ZapConfig, cores ...zapcore.Core) (logger, error) {
	zc, err := zapConfigFromEnv(config.Env)
	if err != nil {
		return nil, err
	}

	if config != nil && config.CustomZapper != nil {
		err = config.CustomZapper.Validate()
		if err != nil {
			return nil, err
		}

		zc.EncoderConfig.FunctionKey = config.CustomZapper.customFunctionKey

		err = zap.RegisterSink(config.CustomZapper.customWriterKey, func(_ *url.URL) (zap.Sink, error) {
			return config.CustomZapper.zw, nil
		})
		if err != nil {
			return nil, err
		}

		zc.OutputPaths = []string{config.CustomZapper.customPath}
	}

	switch config.Env {
	case backend.EnvDevelopment:
		return newZapDevLogger(zc)
	case backend.EnvStaging:
		return switchZapProdLogger(config, zc, cores...)
	case backend.EnvProduction:
		return switchZapProdLogger(config, zc, cores...)
	default:
		return nil, errors.New("invalid environment")
	}
}

func zapConfigFromEnv(env backend.Environment) (zap.Config, error) {
	switch env {
	case backend.EnvDevelopment:
		return zap.NewDevelopmentConfig(), nil
	case backend.EnvStaging:
		fallthrough
	case backend.EnvProduction:
		return zap.NewProductionConfig(), nil
	default:
		return zap.Config{}, errors.New("invalid environment")
	}
}

// ZapDevelopmentLogger is a decorator for zap.Logger that implements
// the logger interface for development environments.
type ZapDevelopmentLogger struct {
	logger  *zap.Logger
	sugared *zap.SugaredLogger
}

func newZapDevLogger(config zap.Config) (*ZapDevelopmentLogger, error) {
	zl, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &ZapDevelopmentLogger{
		logger:  zl,
		sugared: zl.Sugar(),
	}, nil
}

// ZapProductionLogger is a decorator for zap.Logger that implements the
// logger interface for production environments. It supports structured
// logging with zap fields.
type ZapProductionLogger struct {
	logger *zap.Logger
}

func newZapProdLogger(z *zap.Logger) *ZapProductionLogger {
	return &ZapProductionLogger{
		logger: z,
	}
}

// switchZapProdLogger returns a configured production logger based on whether
// telemetry is enabled.
//
// The function signature also takes optional cores. These cores are appended to
// the logger.
func switchZapProdLogger(zc *ZapConfig, zd zap.Config, cores ...zapcore.Core) (logger, error) {
	if zc.WithTelemetry && zc.OtelLoggerProvider == nil {
		return nil, errors.New("telemetry enabled but no logger provider")
	}

	const totalDefaultCores = 2

	allCores := make([]zapcore.Core, 0, len(cores)+totalDefaultCores)
	allCores = append(allCores, cores...)

	if zc.WithTelemetry && zc.OtelLoggerProvider != nil {
		prodCores := []zapcore.Core{
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(os.Stdout),
				zapcore.InfoLevel,
			),
			otelzap.NewCore(zc.OtelServiceName, otelzap.WithLoggerProvider(zc.OtelLoggerProvider)),
		}
		allCores = append(allCores, prodCores...)
	}
	zl, err := zd.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		if len(allCores) > 0 {
			return zapcore.NewTee(allCores...)
		}
		return c
	}))
	if err != nil {
		return nil, err
	}
	return newZapProdLogger(zl), nil
}

// Info logs an info level message using the development logger.
func (l *ZapDevelopmentLogger) Info(msg string, args ...any) {
	if len(args) == 0 {
		l.sugared.Info(msg)
		return
	}
	l.sugared.Infow(msg, args...)
}

// Debug logs a debug level message using the development logger.
func (l *ZapDevelopmentLogger) Debug(msg string, args ...any) {
	if len(args) == 0 {
		l.sugared.Debug(msg)
		return
	}
	l.sugared.Debugw(msg, args...)
}

// Warn logs a warn level message using the development logger.
func (l *ZapDevelopmentLogger) Warn(msg string, args ...any) {
	if len(args) == 0 {
		l.sugared.Warn(msg)
		return
	}
	l.sugared.Warnw(msg, args...)
}

// Error logs an error level message using the development logger.
func (l *ZapDevelopmentLogger) Error(msg string, args ...any) {
	if len(args) == 0 {
		l.sugared.Error(msg)
		return
	}
	l.sugared.Errorw(msg, args...)
}

// Sync flushes the development logger, handling ENOTTY errors gracefully.
func (l *ZapDevelopmentLogger) Sync() error {
	// Without this mess here, Zap will error on any exit. This has something to
	// do with something about file writing. Here's a thread related to
	// this issue:
	// https://github.com/uber-go/zap/issues/880
	//
	// The solution was nabbed from:
	// https://github.com/uber-go/zap/issues/991#issuecomment-962098428

	err := l.logger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		return err
	}
	return nil
}

// zapFieldsFromArgs takes an even number of arguments any and returns a
// slice of the arguments all of the type zap.Field.
//
// There are two cases in which key-value pairs may be missing or may be
// skipped entirely. In the first, If an odd number of arguments are provided,
// the output is truncated to the length of  len(args) - 1. The second case
// is when a key is not a string. When this happens, the key-value pair of the
// non-string key is skipped entirely. It will not appear in the resulting log
// output.
//
// NOTE: Do not use this function directly on a parameter argument. I recommend
// storing the return value of this function in a variable first, THEN using
// that variable for the function argument. Here's an example:
//
// Suppose we have a function func. It has this signature:
// func(args ...any)
//
// Do not do this: func(zapFieldsFromArgs(args)...)
//
// Instead do this:
// fields := zapFieldsFromArgs(args)
// func(fields...)
//
// For some reason, doing the first way causes the fields parameter to have
// nothing in it. It might have something to do with the Go compiler trying
// to optimize an edge case, or it could be something related to expansion
// timing or slice header corruption.
func zapFieldsFromArgs(args ...any) []zap.Field {
	argsLen := len(args)
	if argsLen%2 != 0 {
		argsLen--
		// As an aside, I don't know if this should be an error or not.
	}

	// pairSize describes the number of indexes a key-value pair occupies in
	// an array. The key occupies one index and the value occupies the next
	// consecutive index.
	const pairSize = 2

	fields := make([]zap.Field, 0, argsLen)
	for i := 0; i < argsLen; i += pairSize {
		key, ok := args[i].(string)
		if !ok {
			continue // Skip non-string keys
		}
		fields = append(fields, zap.Any(key, args[i+1]))
	}
	return fields
}

// Info logs an info level message using the production logger with optional structured fields.
func (l *ZapProductionLogger) Info(msg string, args ...any) {
	fields := zapFieldsFromArgs(args...)
	l.logger.Info(msg, fields...)
}

// Debug logs a debug level message using the production logger with optional structured fields.
func (l *ZapProductionLogger) Debug(msg string, args ...any) {
	fields := zapFieldsFromArgs(args...)
	l.logger.Debug(msg, fields...)
}

// Warn logs a warn level message using the production logger with optional structured fields.
func (l *ZapProductionLogger) Warn(msg string, args ...any) {
	fields := zapFieldsFromArgs(args...)
	l.logger.Warn(msg, fields...)
}

// Error logs an error level message using the production logger with optional structured fields.
func (l *ZapProductionLogger) Error(msg string, args ...any) {
	fields := zapFieldsFromArgs(args...)
	l.logger.Error(msg, fields...)
}

// Sync flushes the production logger, handling ENOTTY errors gracefully.
func (l *ZapProductionLogger) Sync() error {
	// Without this mess here, Zap will error on any exit. This has something to
	// do with something about file writing. Here's a thread related to
	// this issue:
	// https://github.com/uber-go/zap/issues/880
	//
	// The solution was nabbed from:
	// https://github.com/uber-go/zap/issues/991#issuecomment-962098428

	err := l.logger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		return err
	}
	return nil
}

// NewPrettySlog creates a logger using the stdlib `slog` package.
// func NewPrettySlog() *slog.Logger {
// 	slogHandlerOptions := &slog.HandlerOptions{
// 		Level: slog.LevelInfo,
// 	}
//
// 	return slog.New(New(slogHandlerOptions))
// }
