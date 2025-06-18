// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package test_test

import (
	"testing"

	"backend.brokedaear.com/internal/common/tests/assert"
	"backend.brokedaear.com/internal/common/tests/test"
)

func TestNewMockLogger(t *testing.T) {
	logger := test.NewMockLogger()
	assert.NotEqual(t, logger, nil)
}

// testMockInterface defines the interface we want to test against.
type testMockLoggerInterface interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Sync() error
}

func TestMockLoggerMethods(t *testing.T) {
	t.Run("has methods", func(t *testing.T) {
		logger := test.NewMockLogger()
		_, ok := logger.(testMockLoggerInterface)
		if !ok {
			t.Error("MockLogger does not implement Logger interface")
		}
		err := logger.Sync()
		assert.NoError(t, err)
	})

	t.Run("methods can be called", func(t *testing.T) {
		// Test that all methods can be called without panicking.

		logger := test.NewMockLogger()

		logger.Info("info message")
		logger.Debug("debug message")
		logger.Warn("warn message")
		logger.Error("error message")

		// Test with arguments.

		logger.Info("info with args", "key1", "value1", "key2", 42)
		logger.Debug("debug with args", "debug", true)
		logger.Warn("warn with args", "count", 100)
		logger.Error("error with args", "error", "test error")

		err := logger.Sync()
		assert.NoError(t, err)
	})
}
