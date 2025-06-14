// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// newLoggerExporter creates a new log exporter.
func newLoggerExporter(ctx context.Context, config ExporterConfig) (log.Exporter, error) {
	switch config.Type {
	case ExporterTypeGRPC:
		return newGRPCExporter(ctx, config)
	case ExporterTypeHTTP:
		return newHTTPExporter(ctx, config)
	case ExporterTypeStdout:
		return newStdoutExporter(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", config.Type)
	}
}

func newGRPCExporter(ctx context.Context, config ExporterConfig) (log.Exporter, error) {
	var opts []otlploggrpc.Option

	if config.Endpoint != "" {
		opts = append(opts, otlploggrpc.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		opts = append(opts, otlploggrpc.WithInsecure())
	}

	if len(config.Headers) > 0 {
		opts = append(opts, otlploggrpc.WithHeaders(config.Headers))
	}

	return otlploggrpc.New(ctx, opts...)
}

func newHTTPExporter(ctx context.Context, config ExporterConfig) (log.Exporter, error) {
	var opts []otlploghttp.Option

	if config.Endpoint != "" {
		opts = append(opts, otlploghttp.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}

	if len(config.Headers) > 0 {
		opts = append(opts, otlploghttp.WithHeaders(config.Headers))
	}

	return otlploghttp.New(ctx, opts...)
}

func newStdoutExporter(_ context.Context, config ExporterConfig) (log.Exporter, error) {
	var opts []stdoutlog.Option

	// Optional: write to a file instead of stdout
	if config.Endpoint != "" {
		file, err := os.OpenFile(config.Endpoint, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", config.Endpoint, err)
		}
		opts = append(opts, stdoutlog.WithWriter(file))
	}

	// Pretty print for better readability (default is false)
	opts = append(opts, stdoutlog.WithPrettyPrint())

	return stdoutlog.New(opts...)
}

// newMetricExporter creates a metric exporter.
func newMetricExporter(ctx context.Context, config ExporterConfig) (metric.Exporter, error) {
	switch config.Type {
	case ExporterTypeGRPC:
		return newGRPCMetricExporter(ctx, config)
	case ExporterTypeHTTP:
		return newHTTPMetricExporter(ctx, config)
	case ExporterTypeStdout:
		return newStdoutMetricExporter(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported metric exporter type: %s", config.Type)
	}
}

func newGRPCMetricExporter(ctx context.Context, config ExporterConfig) (metric.Exporter, error) {
	var opts []otlpmetricgrpc.Option

	if config.Endpoint != "" {
		opts = append(opts, otlpmetricgrpc.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	if len(config.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(config.Headers))
	}

	return otlpmetricgrpc.New(ctx, opts...)
}

func newHTTPMetricExporter(ctx context.Context, config ExporterConfig) (metric.Exporter, error) {
	var opts []otlpmetrichttp.Option

	if config.Endpoint != "" {
		opts = append(opts, otlpmetrichttp.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	if len(config.Headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(config.Headers))
	}

	return otlpmetrichttp.New(ctx, opts...)
}

func newStdoutMetricExporter(_ context.Context, config ExporterConfig) (metric.Exporter, error) {
	var opts []stdoutmetric.Option

	// Optional: write to a file instead of stdout
	if config.Endpoint != "" {
		file, err := os.OpenFile(config.Endpoint, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", config.Endpoint, err)
		}
		opts = append(opts, stdoutmetric.WithWriter(file))
	}

	// Pretty print for better readability
	opts = append(opts, stdoutmetric.WithPrettyPrint())

	return stdoutmetric.New(opts...)
}

// newTraceExporter creates a trace exporter.
func newTraceExporter(ctx context.Context, config ExporterConfig) (trace.SpanExporter, error) {
	switch config.Type {
	case ExporterTypeGRPC:
		return newGRPCTraceExporter(ctx, config)
	case ExporterTypeHTTP:
		return newHTTPTraceExporter(ctx, config)
	case ExporterTypeStdout:
		return newStdoutTraceExporter(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported trace exporter type: %s", config.Type)
	}
}

func newGRPCTraceExporter(ctx context.Context, config ExporterConfig) (trace.SpanExporter, error) {
	var opts []otlptracegrpc.Option

	if config.Endpoint != "" {
		opts = append(opts, otlptracegrpc.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	if len(config.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(config.Headers))
	}

	return otlptracegrpc.New(ctx, opts...)
}

func newHTTPTraceExporter(ctx context.Context, config ExporterConfig) (trace.SpanExporter, error) {
	var opts []otlptracehttp.Option

	if config.Endpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpoint(config.Endpoint))
	}

	if config.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if len(config.Headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(config.Headers))
	}

	return otlptracehttp.New(ctx, opts...)
}

func newStdoutTraceExporter(_ context.Context, config ExporterConfig) (
	trace.SpanExporter,
	error,
) {
	var opts []stdouttrace.Option

	// Optional: write to a file instead of stdout
	if config.Endpoint != "" {
		file, err := os.OpenFile(config.Endpoint, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", config.Endpoint, err)
		}
		opts = append(opts, stdouttrace.WithWriter(file))
	}

	// Pretty print for better readability
	opts = append(opts, stdouttrace.WithPrettyPrint())

	return stdouttrace.New(opts...)
}
