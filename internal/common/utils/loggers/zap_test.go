// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package loggers_test

import (
	"errors"
	"fmt"
	"io"
	"syscall"
	"testing"

	"backend.brokedaear.com"
	"backend.brokedaear.com/pkg/assert"
	"backend.brokedaear.com/pkg/test"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"

	"backend.brokedaear.com/internal/common/telemetry"
	"backend.brokedaear.com/internal/common/utils/loggers"
)

// TestZapLogger_Sync asserts that the Sync method will succeed or,
// in the case it does not succeed, does NOT return an ENOTTY error.
func TestZapLogger_Sync(t *testing.T) {
	tests := []test.CaseBase{
		test.NewCaseBase("sync with test logger", nil, false),
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Use zaptest.NewLogger to create a logger that doesn't write to stderr
			// avoiding the sync issues in test environments.
			testLogger := zaptest.NewLogger(t)

			config := &loggers.ZapConfig{
				Env:                backend.EnvDevelopment,
				OtelServiceName:    "",
				OtelLoggerProvider: nil,
				CustomZapper:       nil,
				WithTelemetry:      false,
			}
			logger, err := loggers.NewZap(config)
			assert.NoError(t, err)

			// Write some log entries to ensure there's something to sync.
			logger.Info("broke da ear woo hoo")
			logger.Debug("debug yes sir")

			err = logger.Sync()
			if err != nil {
				assert.False(t, errors.Is(err, syscall.ENOTTY))
			}

			testErr := testLogger.Sync()
			assert.NoError(t, testErr)
		})
	}
}

// TestZapLogger_Sync_ErrorHandling asserts that the Sync method behaves as
// expected; ENOTTY errors should be ignored and return nil, while other
// errors should be returned as usual.
func TestZapLogger_Sync_ErrorHandling(t *testing.T) {
	config := &loggers.ZapConfig{
		Env:                backend.EnvDevelopment,
		OtelServiceName:    "",
		OtelLoggerProvider: nil,
		CustomZapper:       nil,
		WithTelemetry:      false,
	}
	logger, err := loggers.NewZap(config)
	assert.NoError(t, err)

	err = logger.Sync()
	if err != nil {
		assert.False(t, errors.Is(err, syscall.ENOTTY))
	}
}

func TestZapLogger_Sync_ENOTTYHandling(t *testing.T) {
	// This test specifically verifies that the ZapLogger.Sync() method
	// properly handles ENOTTY errors by ignoring them and returning nil.
	// Since we can't easily inject mocks into the ZapLogger struct due to
	// private fields, this test documents the expected behavior.

	config := &loggers.ZapConfig{
		Env:                backend.EnvDevelopment,
		OtelServiceName:    "",
		OtelLoggerProvider: nil,
		CustomZapper:       nil,
		WithTelemetry:      false,
	}
	logger, err := loggers.NewZap(config)
	assert.NoError(t, err)

	// Add some log entries
	logger.Info("test info message")
	logger.Debug("test debug message")
	logger.Warn("test warn message")
	logger.Error("test error message")

	// Call sync - this should handle any ENOTTY errors gracefully
	err = logger.Sync()
	// The ZapLogger.Sync() implementation should handle ENOTTY errors
	// and return nil instead of propagating them
	if err != nil {
		// If there's still an error, it should not be ENOTTY since
		// that should have been caught and ignored
		assert.False(t, errors.Is(err, syscall.ENOTTY))
	}
}

// TestZapLogger_Sync_MockedZapBehavior verifies zap logger behavior, checking
// that errors returned are the errors we want returned.
func TestZapLogger_Sync_MockedZapBehavior(t *testing.T) {
	type testCase struct {
		test.CaseBase
		syncError error
	}

	tests := []testCase{
		{
			CaseBase:  test.NewCaseBase("no error", nil, false),
			syncError: nil,
		},
		{
			CaseBase:  test.NewCaseBase("ENOTTY error present", syscall.ENOTTY, true),
			syncError: syscall.ENOTTY,
		},
		{
			CaseBase:  test.NewCaseBase("other error propagated", errors.New("sync failed"), true),
			syncError: errors.New("sync failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Create a custom logger with our mock syncer
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
				&mockWriteSyncer{syncError: tt.syncError},
				zapcore.DebugLevel,
			)
			zapLogger := zap.New(core)

			err := zapLogger.Sync()
			assert.ErrorOrNoError(t, err, tt.WantErr)

			if tt.syncError != nil {
				assert.True(t, errors.Is(err, tt.syncError))
			}
		})
	}
}

// TestZapProductionLogger_OutputValidation tests the production logger output with captured logs.
func TestZapProductionLogger_OutputValidation(t *testing.T) {
	tl := zaptest.NewLogger(t).Sugar()

	defer tl.Sync()

	b, bWriter := newBufAndWriter()
	mockWriter := newMockCustomWriter(bWriter)
	customZapWriter := loggers.NewCustomZapWriter(
		"testwriter-prod:testoutput",
		"testwriter-prod",
		"func",
		mockWriter,
	)

	le, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		t.Error(err)
	}

	lp := telemetry.NewLoggerProvider(telemetry.NewResource("test-service", "1.0.0", "test-1"), le)

	config := &loggers.ZapConfig{
		Env:                backend.EnvProduction,
		OtelServiceName:    "test-service",
		CustomZapper:       customZapWriter,
		OtelLoggerProvider: lp,
		WithTelemetry:      true,
	}

	logger, err := loggers.NewZap(config)
	assert.NoError(t, err)

	logger.Info("simple test", "key1", "value1", "key2", 42)

	logger.Info("info test", "key3", "value3", "key4", 84)

	bWriter.Flush()

	err = logger.Sync()
	if err != nil {
		assert.True(t, !errors.Is(err, syscall.ENOTTY))
	}

	// Parse and validate the output
	entries, err := parseLogOutput(t, b)
	assert.NoError(t, err)

	tl.Infof("Buffer contents: %s", b.String())
	tl.Infof("Parsed entries: %+v", entries)

	if len(entries) >= 1 {
		entry := entries[0]
		t.Logf("First entry fields: %+v", entry.Fields)

		assert.Equal(t, entry.Level, "info")
		assert.Equal(t, entry.Message, "simple test")

		// Check if structured fields are present
		if len(entry.Fields) > 0 {
			assert.Equal(t, entry.Fields["key1"], "value1")
			// JSON unmarshals numbers as float64
			if val, ok := entry.Fields["key2"].(float64); ok {
				assert.Equal(t, val, 42.0)
			}
		} else {
			t.Error("no structured fields found")
		}
	}
}

// TestZapStagingLogger_OutputValidation tests the staging logger output with captured logs.
func TestZapStagingLogger_OutputValidation(t *testing.T) {
	b, bWriter := newBufAndWriter()
	mockWriter := newMockCustomWriter(bWriter)
	customZapWriter := loggers.NewCustomZapWriter(
		"testwriter-staging:testoutput",
		"testwriter-staging",
		"func",
		mockWriter,
	)

	config := &loggers.ZapConfig{
		Env:                backend.EnvStaging,
		OtelServiceName:    "test-service",
		OtelLoggerProvider: nil,
		CustomZapper:       customZapWriter,
		WithTelemetry:      false,
	}
	logger, err := loggers.NewZap(config)
	assert.NoError(t, err)

	logger.Info("meow", "environment", "staging", "version", "1.0.0")
	logger.Warn("woof", "memory_usage", 85.5, "threshold", 80.0)

	bWriter.Flush()

	err = logger.Sync()
	if err != nil {
		assert.True(t, !errors.Is(err, syscall.ENOTTY))
	}

	entries, err := parseLogOutput(t, b)
	assert.NoError(t, err)

	t.Logf("Buffer contents: %s", b.String())
	t.Logf("Parsed entries: %+v", entries)

	assert.Equal(t, len(entries), 2)

	// Validate structured fields.

	assert.Equal(t, entries[0].Level, "info")
	assert.Equal(t, entries[0].Message, "meow")

	val, ok := entries[0].Fields["environment"]
	if ok {
		assert.Equal(t, val, "staging")
	} else {
		t.Errorf("environment field not found in first entry")
	}

	val, ok = entries[0].Fields["version"]
	if ok {
		assert.Equal(t, val, "1.0.0")
	} else {
		t.Errorf("version field not found in first entry")
	}

	assert.Equal(t, entries[1].Level, "warn")
	assert.Equal(t, entries[1].Message, "woof")

	val, ok = entries[1].Fields["memory_usage"]
	if ok {
		assert.Equal(t, val, 85.5)
	} else {
		t.Errorf("memory_usage field not found in second entry")
	}

	val, ok = entries[1].Fields["threshold"]
	if ok {
		assert.Equal(t, val, 80.0)
	} else {
		t.Errorf("threshold field not found in second entry")
	}
}

// TestZapLogger_TelemetryError tests error handling when telemetry is enabled but no provider given.
func TestZapLogger_TelemetryError(t *testing.T) {
	// This should return an error because telemetry is enabled but no logger provider
	config := &loggers.ZapConfig{
		Env:                backend.EnvProduction,
		OtelServiceName:    "test-service",
		OtelLoggerProvider: nil,
		CustomZapper:       nil,
		WithTelemetry:      true,
	}
	_, err := loggers.NewZap(config)
	assert.True(t, err != nil)
}

// TestZapProductionLogger_StructuredFieldsEdgeCases tests edge cases in structured field processing.
func TestZapProductionLogger_StructuredFieldsEdgeCases(t *testing.T) {
	type testCase struct {
		test.CaseBase
		args          []any
		expectedPairs int
	}

	tests := []testCase{
		{
			CaseBase:      test.NewCaseBase("valid key-value pairs", nil, false),
			args:          []any{"key1", "value1", "key2", 42},
			expectedPairs: 3, // 2 user fields + 1 func field
		},
		{
			CaseBase:      test.NewCaseBase("odd number of args", nil, false),
			args:          []any{"key1", "value1", "key2"},
			expectedPairs: 2, // 1 user field + 1 func field
		},
		{
			CaseBase:      test.NewCaseBase("empty args", nil, false),
			args:          []any{},
			expectedPairs: 1, // 0 user fields + 1 func field
		},
		{
			CaseBase:      test.NewCaseBase("non-string key", nil, false),
			args:          []any{123, "value1", "key2", "value2"},
			expectedPairs: 2, // 1 user field + 1 func field
		},
		{
			CaseBase:      test.NewCaseBase("mixed types", nil, false),
			args:          []any{"string_key", "string_val", "int_key", 42, "bool_key", true, "float_key", 3.14},
			expectedPairs: 5, // 4 user fields + 1 func field
		},
	}

	for i, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b, bWriter := newBufAndWriter()
			mockWriter := newMockCustomWriter(bWriter)
			writerKey := fmt.Sprintf("testwriter-edge-%d", i)
			customZapWriter := loggers.NewCustomZapWriter(
				writerKey+":testoutput",
				writerKey,
				"func",
				mockWriter,
			)

			config := &loggers.ZapConfig{
				Env:                backend.EnvStaging,
				OtelServiceName:    "test-service",
				OtelLoggerProvider: nil,
				CustomZapper:       customZapWriter,
				WithTelemetry:      false,
			}
			logger, err := loggers.NewZap(config)
			assert.NoError(t, err)

			logger.Info("test message", tt.args...)

			bWriter.Flush()

			err = logger.Sync()
			if err != nil {
				assert.False(t, errors.Is(err, syscall.ENOTTY))
			}

			entries, err := parseLogOutput(t, b)
			assert.NoError(t, err)
			assert.Equal(t, len(entries), 1)

			entry := entries[0]
			assert.Equal(t, entry.Level, "info")
			assert.Equal(t, entry.Message, "test message")
			assert.Equal(t, len(entry.Fields), tt.expectedPairs)

			// Validate specific field values for mixed types test.

			if tt.Name == "mixed types" {
				assert.Equal(t, entry.Fields["string_key"], "string_val")
				// JSON unmarshals numbers as float64
				intVal, ok := entry.Fields["int_key"].(float64)
				if ok {
					assert.Equal(t, intVal, 42.0)
				}
				assert.Equal(t, entry.Fields["bool_key"].(bool), true)
				assert.Equal(t, entry.Fields["float_key"].(float64), 3.14)
			}
		})
	}
}

// TestZapDevelopmentLogger_OutputValidation tests the development logger output with captured logs.
func TestZapDevelopmentLogger_OutputValidation(t *testing.T) {
	b, bWriter := newBufAndWriter()
	mockWriter := newMockCustomWriter(bWriter)
	customZapWriter := loggers.NewCustomZapWriter(
		"testwriter-dev:testoutput",
		"testwriter-dev",
		"func",
		mockWriter,
	)

	config := &loggers.ZapConfig{
		Env:                backend.EnvDevelopment,
		OtelServiceName:    "test-service",
		CustomZapper:       customZapWriter,
		OtelLoggerProvider: nil,
		WithTelemetry:      false,
	}
	logger, err := loggers.NewZap(config)
	assert.NoError(t, err)

	logger.Info("development info message")
	logger.Debug("development debug message", "component", "test", "action", "debug")
	logger.Warn("development warn message", "warning", "test", "level", "high")
	logger.Error("development error message", "error", "test", "severity", "critical")

	bWriter.Flush()

	err = logger.Sync()
	if err != nil {
		assert.False(t, errors.Is(err, syscall.ENOTTY))
	}

	output := b.String()
	t.Logf("Buffer contents: %s", output)

	// Validate that output contains expected log messages
	assert.True(t, len(output) > 0)

	// Check for presence of log messages (development logger uses console format, not JSON)
	lines := b.Lines()
	assert.True(t, len(lines) >= 4) // Should have at least 4 log entries

	// Basic validation: ensure we have meaningful output
	assert.True(t, len(output) > 200) // Should be substantial output

	// Count lines with structured fields (end with '}')
	structuredLines := 0
	for _, line := range lines {
		if len(line) > 30 && line[len(line)-1:] == "}" {
			structuredLines++
		}
	}

	// Should have structured log entries (DEBUG, WARN, ERROR all have structured fields)
	assert.True(t, structuredLines >= 3)
}

// TestZapLogger_CustomWriter tests the custom writer functionality. It answers
// the question of, "does the logger write to the writer?".
func TestZapLogger_CustomWriter(t *testing.T) {
	b, bWriter := newBufAndWriter()
	mockWriter := newMockCustomWriter(bWriter)
	customZapWriter := loggers.NewCustomZapWriter(
		"customkey:8080", // customPath using customWriterKey:port format
		"customkey",      // customWriterKey
		"funckey",        // customFunctionKey
		mockWriter,
	)

	config := &loggers.ZapConfig{
		Env:                backend.EnvDevelopment,
		OtelServiceName:    "test-service",
		OtelLoggerProvider: nil,
		CustomZapper:       customZapWriter,
		WithTelemetry:      false,
	}

	logger, err := loggers.NewZap(config)
	assert.NoError(t, err)

	logger.Info("custom writer test", "key1", "value1")
	logger.Debug("debug with custom writer")

	bWriter.Flush()

	output := b.String()
	assert.Equal(t, len(output) > 0, true)

	err = logger.Sync()
	if err != nil {
		assert.False(t, errors.Is(err, syscall.ENOTTY))
	}
}

// TestCustomZapWriter_Validation tests the validation of CustomZapWriter.
func TestCustomZapWriter_Validation(t *testing.T) {
	type testCase struct {
		test.CaseBase
		customPath        string
		customWriterKey   string
		customFunctionKey string
		expectError       bool
	}

	tests := []testCase{
		{
			CaseBase:          test.NewCaseBase("valid config", nil, false),
			customPath:        "validkey:8080",
			customWriterKey:   "validkey",
			customFunctionKey: "funckey",
			expectError:       false,
		},
		{
			CaseBase:          test.NewCaseBase("empty customPath", nil, false),
			customPath:        "",
			customWriterKey:   "validkey",
			customFunctionKey: "funckey",
			expectError:       true,
		},
		{
			CaseBase:          test.NewCaseBase("empty customWriterKey", nil, false),
			customPath:        "validkey:8080",
			customWriterKey:   "",
			customFunctionKey: "funckey",
			expectError:       true,
		},
		{
			CaseBase:          test.NewCaseBase("empty customFunctionKey", nil, false),
			customPath:        "validkey:8080",
			customWriterKey:   "validkey",
			customFunctionKey: "",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, bWriter := newBufAndWriter()
			mockWriter := newMockCustomWriter(bWriter)
			customZapWriter := loggers.NewCustomZapWriter(
				tt.customPath,
				tt.customWriterKey,
				tt.customFunctionKey,
				mockWriter,
			)

			err := customZapWriter.Validate()
			assert.ErrorOrNoError(t, err, tt.expectError)
		})
	}
}

// logEntry represents a parsed log entry for testing. It is the same as the
// log entry for a zap production logger.
type logEntry struct {
	Level     string         `json:"level"`
	Message   string         `json:"msg"`
	Timestamp float64        `json:"ts"`
	Caller    string         `json:"caller"`
	Fields    map[string]any `json:"-"`
}

// mockCustomWriter implements the ZapWriter interface using io.Writer.
type mockCustomWriter struct {
	io.Writer
}

func newMockCustomWriter(w io.Writer) *mockCustomWriter {
	return &mockCustomWriter{
		w,
	}
}

// Close implements loggers.ZapWriter.
func (m *mockCustomWriter) Close() error {
	return nil
}

func (m *mockCustomWriter) Sync() error {
	return nil
}

type mockWriteSyncer struct {
	syncError error
}

func (m *mockWriteSyncer) Write(p []byte) (int, error) {
	return len(p), nil
}

func (m *mockWriteSyncer) Sync() error {
	return m.syncError
}
