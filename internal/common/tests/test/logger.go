// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package test

// MockLogger mocks the Logger interface.
type MockLogger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Sync() error
}

// mockLogger implements the MockLogger interface.
type mockLogger struct {
	logs []string
}

func NewMockLogger() MockLogger {
	return &mockLogger{
		logs: make([]string, 0),
	}
}

func (m *mockLogger) Info(msg string, _ ...any) {
	m.logs = append(m.logs, "INFO: "+msg)
}

func (m *mockLogger) Debug(msg string, _ ...any) {
	m.logs = append(m.logs, "DEBUG: "+msg)
}

func (m *mockLogger) Warn(msg string, _ ...any) {
	m.logs = append(m.logs, "WARN: "+msg)
}

func (m *mockLogger) Error(msg string, _ ...any) {
	m.logs = append(m.logs, "ERROR: "+msg)
}

func (m *mockLogger) Sync() error {
	return nil
}
