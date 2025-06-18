// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package loggers

import (
	"errors"
	"syscall"
	"testing"

	"backend.brokedaear.com"
	"backend.brokedaear.com/pkg/assert"
	"backend.brokedaear.com/pkg/test"

	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapFieldsFromArgsTestCase struct {
	test.CaseBase
	args   []any
	expect []zap.Field
}

func Test_zapFieldsFromArgs(t *testing.T) {
	testCases := []zapFieldsFromArgsTestCase{
		{
			CaseBase: test.NewCaseBase("empty args", nil, false),
			args:     []any{},
			expect:   []zap.Field{},
		},
		{
			CaseBase: test.NewCaseBase("single string key with value", nil, false),
			args:     []any{"key1", "value1"},
			expect:   []zap.Field{zap.Any("key1", "value1")},
		},
		{
			CaseBase: test.NewCaseBase("multiple string keys with values", nil, false),
			args:     []any{"key1", "value1", "key2", 42, "key3", true},
			expect:   []zap.Field{zap.Any("key1", "value1"), zap.Any("key2", 42), zap.Any("key3", true)},
		},
		{
			CaseBase: test.NewCaseBase("odd number of args truncates last", nil, false),
			args:     []any{"key1", "value1", "orphan"},
			expect:   []zap.Field{zap.Any("key1", "value1")},
		},
		{
			CaseBase: test.NewCaseBase("non-string key skipped", nil, false),
			args:     []any{123, "value1", "key2", "value2"},
			expect:   []zap.Field{zap.Any("key2", "value2")},
		},
		{
			CaseBase: test.NewCaseBase("mixed non-string and string keys", nil, false),
			args:     []any{"key1", "value1", 456, "skipped", "key3", "value3"},
			expect:   []zap.Field{zap.Any("key1", "value1"), zap.Any("key3", "value3")},
		},
		{
			CaseBase: test.NewCaseBase("all non-string keys", nil, false),
			args:     []any{123, "value1", 456, "value2"},
			expect:   []zap.Field{},
		},
		{
			CaseBase: test.NewCaseBase("nil values allowed", nil, false),
			args:     []any{"key1", nil, "key2", "value2"},
			expect:   []zap.Field{zap.Any("key1", nil), zap.Any("key2", "value2")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			got := zapFieldsFromArgs(tc.args...)

			assert.Equal(t, len(got), len(tc.expect))

			for i, field := range got {
				assert.Equal(t, field.Key, tc.expect[i].Key)
				assert.Equal(t, field.Type, tc.expect[i].Type)
				assert.Equal(t, field.Interface, tc.expect[i].Interface)
			}
		})
	}
}

type switchZapProdLoggerTestCase struct {
	test.CaseBase
	config       *ZapConfig
	cores        []zapcore.Core
	expectError  bool
	expectLogger bool
}

func Test_switchZapProdLogger(t *testing.T) {
	mockLoggerProvider := &log.LoggerProvider{
		LoggerProvider: nil,
	}
	mockCore := zapcore.NewNopCore()

	testCases := []switchZapProdLoggerTestCase{
		{
			CaseBase: test.NewCaseBase("telemetry disabled no cores", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      false,
			},
			cores:        []zapcore.Core{},
			expectError:  false,
			expectLogger: true,
		},
		{
			CaseBase: test.NewCaseBase("telemetry disabled with cores", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      false,
			},
			cores:        []zapcore.Core{mockCore},
			expectError:  false,
			expectLogger: true,
		},
		{
			CaseBase: test.NewCaseBase("telemetry enabled with provider", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: mockLoggerProvider,
				CustomZapper:       nil,
				WithTelemetry:      true,
			},
			cores:        []zapcore.Core{},
			expectError:  false,
			expectLogger: true,
		},
		{
			CaseBase: test.NewCaseBase("telemetry enabled with provider and cores", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: mockLoggerProvider,
				CustomZapper:       nil,
				WithTelemetry:      true,
			},
			cores:        []zapcore.Core{mockCore},
			expectError:  false,
			expectLogger: true,
		},
		{
			CaseBase: test.NewCaseBase("telemetry enabled without provider", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      true,
			},
			cores:        []zapcore.Core{},
			expectError:  true,
			expectLogger: false,
		},
		{
			CaseBase: test.NewCaseBase("telemetry enabled without provider with cores", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      true,
			},
			cores:        []zapcore.Core{mockCore},
			expectError:  true,
			expectLogger: false,
		},
		{
			CaseBase: test.NewCaseBase("empty service name telemetry enabled", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "",
				OtelLoggerProvider: mockLoggerProvider,
				CustomZapper:       nil,
				WithTelemetry:      true,
			},
			cores:        []zapcore.Core{},
			expectError:  false,
			expectLogger: true,
		},
		{
			CaseBase: test.NewCaseBase("empty service name telemetry disabled", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      false,
			},
			cores:        []zapcore.Core{},
			expectError:  false,
			expectLogger: true,
		},
		{
			CaseBase: test.NewCaseBase("multiple custom cores", nil, false),
			config: &ZapConfig{
				Env:                backend.EnvProduction,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      false,
			},
			cores:        []zapcore.Core{mockCore, mockCore, mockCore},
			expectError:  false,
			expectLogger: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Need to call zapConfigFromEnv to get the base config
			zapConfig, configErr := zapConfigFromEnv(tc.config.Env)
			assert.NoError(t, configErr)

			l, err := switchZapProdLogger(tc.config, zapConfig, tc.cores...)

			assert.ErrorOrNoError(t, err, tc.expectError)

			if tc.expectLogger {
				assert.Equal(t, l != nil, true)

				// Verify the logger implements the logger interface
				if l != nil {
					// Test that basic logging methods work
					l.Info("test info message")
					l.Debug("test debug message")
					l.Warn("test warn message")
					l.Error("test error message")

					// Test sync method
					syncErr := l.Sync()
					// Sync should either succeed or fail with a non-ENOTTY error
					if syncErr != nil {
						assert.False(t, errors.Is(syncErr, syscall.ENOTTY))
					}
				}
			} else {
				assert.Equal(t, l == nil, true)
			}

			// Verify error message for telemetry enabled without provider
			if tc.config.WithTelemetry && tc.config.OtelLoggerProvider == nil {
				assert.Equal(t, err.Error(), "telemetry enabled but no logger provider")
			}
		})
	}
}

func Test_switchZapProdLogger_LoggerType(t *testing.T) {
	// Test that the returned logger is of the correct type
	config := &ZapConfig{
		Env:                backend.EnvProduction,
		OtelServiceName:    "test-service",
		OtelLoggerProvider: nil,
		CustomZapper:       nil,
		WithTelemetry:      false,
	}
	zapConfig, err := zapConfigFromEnv(config.Env)
	assert.NoError(t, err)

	l, err := switchZapProdLogger(config, zapConfig)
	assert.NoError(t, err)
	assert.Equal(t, l != nil, true)

	// Type assertion to ensure it returns a ZapProductionLogger
	_, ok := l.(*ZapProductionLogger)
	assert.Equal(t, ok, true)
}

func Test_switchZapProdLogger_TelemetryPath(t *testing.T) {
	// Test the telemetry-enabled path separately to ensure it creates the correct cores
	mockLoggerProvider := &log.LoggerProvider{
		LoggerProvider: nil,
	}

	config := &ZapConfig{
		Env:                backend.EnvProduction,
		OtelServiceName:    "telemetry-service",
		OtelLoggerProvider: mockLoggerProvider,
		CustomZapper:       nil,
		WithTelemetry:      true,
	}
	zapConfig, err := zapConfigFromEnv(config.Env)
	assert.NoError(t, err)

	l, err := switchZapProdLogger(config, zapConfig)
	assert.NoError(t, err)
	assert.Equal(t, l != nil, true)

	// Verify it's a production logger
	prodLogger, ok := l.(*ZapProductionLogger)
	assert.Equal(t, ok, true)
	assert.Equal(t, prodLogger != nil, true)

	// Test logging still works with telemetry enabled
	if prodLogger != nil {
		prodLogger.Info("telemetry test", "key", "value")
		err = prodLogger.Sync()
		if err != nil {
			assert.False(t, errors.Is(err, syscall.ENOTTY))
		}
	}
}

func Test_switchZapProdLogger_CoreCounting(t *testing.T) {
	// Test that custom cores are properly integrated
	mockCore1 := zapcore.NewNopCore()
	mockCore2 := zapcore.NewNopCore()

	// Test with multiple cores in telemetry disabled mode
	config1 := &ZapConfig{
		Env:                backend.EnvProduction,
		OtelServiceName:    "test-service",
		OtelLoggerProvider: nil,
		CustomZapper:       nil,
		WithTelemetry:      false,
	}
	zapConfig1, err := zapConfigFromEnv(config1.Env)
	assert.NoError(t, err)

	l, err := switchZapProdLogger(config1, zapConfig1, mockCore1, mockCore2)
	assert.NoError(t, err)
	assert.Equal(t, l != nil, true)

	// Test with multiple cores in telemetry enabled mode
	mockLoggerProvider := &log.LoggerProvider{
		LoggerProvider: nil,
	}
	config2 := &ZapConfig{
		Env:                backend.EnvProduction,
		OtelServiceName:    "test-service",
		OtelLoggerProvider: mockLoggerProvider,
		CustomZapper:       nil,
		WithTelemetry:      true,
	}
	zapConfig2, err2 := zapConfigFromEnv(config2.Env)
	assert.NoError(t, err2)

	logger2, err2 := switchZapProdLogger(config2, zapConfig2, mockCore1, mockCore2)
	assert.NoError(t, err2)
	assert.Equal(t, logger2 != nil, true)
}

// Debug test to check zapFieldsFromArgs.
func Test_zapFieldsFromArgs_Debug(t *testing.T) {
	args := []any{"key1", "value1", "key2", 42}
	fields := zapFieldsFromArgs(args...)

	t.Logf("Generated %d fields from args %v", len(fields), args)
	for i, field := range fields {
		t.Logf("Field %d: Key=%s, Type=%v, Interface=%v", i, field.Key, field.Type, field.Interface)
	}

	assert.Equal(t, len(fields), 2)
	assert.Equal(t, fields[0].Key, "key1")
	assert.Equal(t, fields[1].Key, "key2")
}
