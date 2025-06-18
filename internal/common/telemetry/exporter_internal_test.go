// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"os"
	"path/filepath"
	"testing"

	"backend.brokedaear.com/pkg/assert"
	"backend.brokedaear.com/pkg/test"
)

type LoggerExporterTestCase struct {
	test.CaseBase
	Config ExporterConfig
}

func TestNewLoggerExporter(t *testing.T) {
	ctx := t.Context()

	tests := []LoggerExporterTestCase{
		{
			CaseBase: test.NewCaseBase("grpc exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "localhost:4317",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("http exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeHTTP,
				Endpoint: "localhost:4318",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("stdout exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: "",
				Insecure: false,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("stdout exporter with file", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: filepath.Join(t.TempDir(), "test-logs.json"),
				Insecure: false,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"unsupported exporter type",
				"unsupported exporter type",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterType(99),
				Endpoint: "",
				Insecure: false,
				Headers:  make(map[string]string),
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				exporter, err := newLoggerExporter(ctx, tt.Config)
				assert.ErrorOrNoError(t, err, tt.WantErr)
				if !tt.WantErr {
					if exporter == nil {
						t.Error("expected non-nil exporter")
					} else {
						defer func() {
							_ = exporter.Shutdown(ctx)
						}()
					}
				}
			},
		)
	}
}

type MetricExporterTestCase struct {
	test.CaseBase
	Config ExporterConfig
}

func TestNewMetricExporter(t *testing.T) {
	ctx := t.Context()

	tests := []MetricExporterTestCase{
		{
			CaseBase: test.NewCaseBase("grpc metric exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "localhost:4317",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("http metric exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeHTTP,
				Endpoint: "localhost:4318",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("stdout metric exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: "",
				Insecure: false,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("stdout metric exporter with file", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: filepath.Join(t.TempDir(), "test-metrics.json"),
				Insecure: false,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"unsupported metric exporter type",
				"unsupported metric exporter type",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterType(99),
				Endpoint: "",
				Insecure: false,
				Headers:  make(map[string]string),
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				exporter, err := newMetricExporter(ctx, tt.Config)
				assert.ErrorOrNoError(t, err, tt.WantErr)
				if !tt.WantErr {
					if exporter == nil {
						t.Error("expected non-nil exporter")
					} else {
						defer func() {
							_ = exporter.Shutdown(ctx)
						}()
					}
				}
			},
		)
	}
}

type TraceExporterTestCase struct {
	test.CaseBase
	Config ExporterConfig
}

func TestNewTraceExporter(t *testing.T) {
	ctx := t.Context()

	tests := []TraceExporterTestCase{
		{
			CaseBase: test.NewCaseBase("grpc trace exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "localhost:4317",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("http trace exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeHTTP,
				Endpoint: "localhost:4318",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("stdout trace exporter", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: "",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase("stdout trace exporter with file", nil, false),
			Config: ExporterConfig{
				Type:     ExporterTypeStdout,
				Endpoint: filepath.Join(t.TempDir(), "test-traces.json"),
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
		{
			CaseBase: test.NewCaseBase(
				"unsupported trace exporter type",
				"unsupported trace exporter type",
				true,
			),
			Config: ExporterConfig{
				Type:     ExporterType(99),
				Endpoint: "",
				Insecure: true,
				Headers:  make(map[string]string),
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				exporter, err := newTraceExporter(ctx, tt.Config)
				assert.ErrorOrNoError(t, err, tt.WantErr)
				if !tt.WantErr {
					if exporter == nil {
						t.Error("expected non-nil exporter")
					} else {
						defer func() {
							_ = exporter.Shutdown(ctx)
						}()
					}
				}
			},
		)
	}
}

func TestNewStdoutExporterWithInvalidFile(t *testing.T) {
	ctx := t.Context()
	config := ExporterConfig{
		Type:     ExporterTypeStdout,
		Endpoint: "/invalid/path/that/does/not/exist/test.json",
		Insecure: true,
		Headers:  make(map[string]string),
	}

	exporter, err := newStdoutExporter(ctx, config)
	if err == nil {
		t.Error("expected error for invalid file path")
	}
	if exporter != nil {
		t.Error("expected nil exporter on error")
	}
}

func TestNewStdoutMetricExporterWithInvalidFile(t *testing.T) {
	ctx := t.Context()
	config := ExporterConfig{
		Type:     ExporterTypeStdout,
		Endpoint: "/invalid/path/that/does/not/exist/test.json",
		Insecure: true,
		Headers:  make(map[string]string),
	}

	exporter, err := newStdoutMetricExporter(ctx, config)
	if err == nil {
		t.Error("expected error for invalid file path")
	}
	if exporter != nil {
		t.Error("expected nil exporter on error")
	}
}

func TestNewStdoutTraceExporterWithInvalidFile(t *testing.T) {
	ctx := t.Context()
	config := ExporterConfig{
		Type:     ExporterTypeStdout,
		Endpoint: "/invalid/path/that/does/not/exist/test.json",
		Insecure: true,
		Headers:  make(map[string]string),
	}

	exporter, err := newStdoutTraceExporter(ctx, config)
	if err == nil {
		t.Error("expected error for invalid file path")
	}
	if exporter != nil {
		t.Error("expected nil exporter on error")
	}
}

func TestNewGRPCExporterWithHeaders(t *testing.T) {
	ctx := t.Context()
	config := ExporterConfig{
		Type:     ExporterTypeGRPC,
		Endpoint: "localhost:4317",
		Insecure: true,
		Headers: map[string]string{
			"api-key":    "test-key",
			"user-agent": "test-agent",
		},
	}

	exporter, err := newGRPCExporter(ctx, config)
	assert.NoError(t, err)
	if exporter == nil {
		t.Error("expected non-nil exporter")
	}
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()
}

func TestNewHTTPExporterWithHeaders(t *testing.T) {
	ctx := t.Context()
	config := ExporterConfig{
		Type:     ExporterTypeHTTP,
		Endpoint: "localhost:4318",
		Insecure: true,
		Headers: map[string]string{
			"api-key":    "test-key",
			"user-agent": "test-agent",
		},
	}

	exporter, err := newHTTPExporter(ctx, config)
	assert.NoError(t, err)
	if exporter == nil {
		t.Error("expected non-nil exporter")
	}
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()
}

func TestStdoutExporterWritesToFile(t *testing.T) {
	ctx := t.Context()
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "test-output.json")

	config := ExporterConfig{
		Type:     ExporterTypeStdout,
		Endpoint: outputFile,
		Insecure: true,
		Headers:  make(map[string]string),
	}

	exporter, err := newStdoutExporter(ctx, config)
	assert.NoError(t, err)
	if exporter == nil {
		t.Error("expected non-nil exporter")
	}
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	// Verify file was created
	_, err = os.Stat(outputFile)
	assert.NoError(t, err)
}

func BenchmarkNewLoggerExporter(b *testing.B) {
	ctx := b.Context()
	config := ExporterConfig{
		Type:     ExporterTypeGRPC,
		Endpoint: "localhost:4317",
		Insecure: true,
		Headers:  make(map[string]string),
	}

	b.ResetTimer()
	for b.Loop() {
		exporter, err := newLoggerExporter(ctx, config)
		if err != nil {
			b.Fatal(err)
		}
		_ = exporter.Shutdown(ctx)
	}
}

func BenchmarkNewMetricExporter(b *testing.B) {
	ctx := b.Context()
	config := ExporterConfig{
		Type:     ExporterTypeGRPC,
		Endpoint: "localhost:4317",
		Insecure: true,
		Headers:  make(map[string]string),
	}

	b.ResetTimer()
	for b.Loop() {
		exporter, err := newMetricExporter(ctx, config)
		if err != nil {
			b.Fatal(err)
		}
		_ = exporter.Shutdown(ctx)
	}
}

func BenchmarkNewTraceExporter(b *testing.B) {
	ctx := b.Context()
	config := ExporterConfig{
		Type:     ExporterTypeGRPC,
		Endpoint: "localhost:4317",
		Insecure: true,
		Headers:  make(map[string]string),
	}

	b.ResetTimer()
	for b.Loop() {
		exporter, err := newTraceExporter(ctx, config)
		if err != nil {
			b.Fatal(err)
		}
		_ = exporter.Shutdown(ctx)
	}
}
