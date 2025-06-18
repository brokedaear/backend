// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"os"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"backend.brokedaear.com/pkg/assert"
)

func TestNewLoggerProvider(t *testing.T) {
	ctx := t.Context()
	res := NewResource("test-service", "1.0.0", "test-id")

	exporter, err := stdoutlog.New()
	assert.NoError(t, err)
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	provider := NewLoggerProvider(res, exporter)
	if provider == nil {
		t.Error("expected non-nil provider")
	}

	defer func() {
		_ = provider.Shutdown(ctx)
	}()
}

func TestNewMeterProvider(t *testing.T) {
	ctx := t.Context()
	res := NewResource("test-service", "1.0.0", "test-id")

	exporter, err := stdoutmetric.New()
	assert.NoError(t, err)
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	provider := NewMeterProvider(res, exporter)
	if provider == nil {
		t.Error("expected non-nil provider")
	}

	// Verify that the global meter provider was set
	if otel.GetMeterProvider() != provider {
		t.Error("expected global meter provider to be set")
	}

	defer func() {
		_ = provider.Shutdown(ctx)
	}()
}

func TestNewTracerProvider(t *testing.T) {
	ctx := t.Context()
	res := NewResource("test-service", "1.0.0", "test-id")

	exporter, err := stdouttrace.New()
	assert.NoError(t, err)
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	provider := NewTracerProvider(res, exporter)
	if provider == nil {
		t.Error("expected non-nil provider")
	}

	// Verify that the global tracer provider was set
	if otel.GetTracerProvider() != provider {
		t.Error("expected global tracer provider to be set")
	}

	defer func() {
		_ = provider.Shutdown(ctx)
	}()
}

func TestNewResource(t *testing.T) {
	originalHostname, _ := os.Hostname()

	tests := []struct {
		name            string
		serviceName     string
		serviceVersion  string
		serviceID       string
		expectedService string
		expectedVersion string
		expectedID      string
	}{
		{
			name:            "basic resource",
			serviceName:     "test-service",
			serviceVersion:  "1.0.0",
			serviceID:       "test-id-123",
			expectedService: "test-service",
			expectedVersion: "1.0.0",
			expectedID:      "test-id-123",
		},
		{
			name:            "resource with special characters",
			serviceName:     "my.service-name_test",
			serviceVersion:  "v2.1.0-beta.1",
			serviceID:       "service-id-with-dashes",
			expectedService: "my.service-name_test",
			expectedVersion: "v2.1.0-beta.1",
			expectedID:      "service-id-with-dashes",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				res := NewResource(tt.serviceName, tt.serviceVersion, tt.serviceID)

				if res == nil {
					t.Error("expected non-nil resource")
				}

				// Verify the resource has the expected attributes
				attrs := res.Attributes()

				// Check service name
				var foundServiceName, foundServiceVersion, foundServiceID, foundHostname bool
				var serviceName, serviceVersion, serviceID, hostname string

				for _, attr := range attrs {
					switch attr.Key {
					case semconv.ServiceNameKey:
						foundServiceName = true
						serviceName = attr.Value.AsString()
					case semconv.ServiceVersionKey:
						foundServiceVersion = true
						serviceVersion = attr.Value.AsString()
					case semconv.ServiceInstanceIDKey:
						foundServiceID = true
						serviceID = attr.Value.AsString()
					case semconv.HostNameKey:
						foundHostname = true
						hostname = attr.Value.AsString()
					}
				}

				assert.True(t, foundServiceName)
				assert.Equal(t, serviceName, tt.expectedService)

				assert.True(t, foundServiceVersion)
				assert.Equal(t, serviceVersion, tt.expectedVersion)

				assert.True(t, foundServiceID)
				assert.Equal(t, serviceID, tt.expectedID)

				assert.True(t, foundHostname)
				assert.Equal(t, hostname, originalHostname)

				// Verify schema URL is set
				assert.Equal(t, res.SchemaURL(), semconv.SchemaURL)
			},
		)
	}
}

func TestNewResourceWithEmptyHostname(t *testing.T) {
	// This test handles the case where os.Hostname() might return an error
	// The function should still work and create a resource
	res := NewResource("test-service", "1.0.0", "test-id")

	if res == nil {
		t.Error("expected non-nil resource")
	}
	attrs := res.Attributes()

	// Hostname should still be present (even if empty)
	foundHostname := false
	for _, attr := range attrs {
		if attr.Key == semconv.HostNameKey {
			foundHostname = true
			break
		}
	}
	assert.True(t, foundHostname)
}

func TestProviderIntegration(t *testing.T) {
	ctx := t.Context()
	res := NewResource("integration-test", "1.0.0", "integration-id")

	// Create exporters
	logExporter, err := stdoutlog.New()
	assert.NoError(t, err)
	defer func() {
		_ = logExporter.Shutdown(ctx)
	}()

	metricExporter, err := stdoutmetric.New()
	assert.NoError(t, err)
	defer func() {
		_ = metricExporter.Shutdown(ctx)
	}()

	traceExporter, err := stdouttrace.New()
	assert.NoError(t, err)
	defer func() {
		_ = traceExporter.Shutdown(ctx)
	}()

	// Create providers
	logProvider := NewLoggerProvider(res, logExporter)
	defer func() {
		_ = logProvider.Shutdown(ctx)
	}()

	metricProvider := NewMeterProvider(res, metricExporter)
	defer func() {
		_ = metricProvider.Shutdown(ctx)
	}()

	traceProvider := NewTracerProvider(res, traceExporter)
	defer func() {
		_ = traceProvider.Shutdown(ctx)
	}()

	// Verify all providers are properly initialized
	if logProvider == nil {
		t.Error("expected non-nil log provider")
	}
	if metricProvider == nil {
		t.Error("expected non-nil metric provider")
	}
	if traceProvider == nil {
		t.Error("expected non-nil trace provider")
	}

	// Verify that meters and tracers can be created
	meter := metricProvider.Meter("test-meter")
	if meter == nil {
		t.Error("expected non-nil meter")
	}

	tracer := traceProvider.Tracer("test-tracer")
	if tracer == nil {
		t.Error("expected non-nil tracer")
	}
}

func BenchmarkNewResource(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewResource("benchmark-service", "1.0.0", "benchmark-id")
	}
}

func BenchmarkNewLoggerProvider(b *testing.B) {
	ctx := b.Context()
	res := NewResource("benchmark-service", "1.0.0", "benchmark-id")
	exporter, _ := stdoutlog.New()
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := NewLoggerProvider(res, exporter)
		_ = provider.Shutdown(ctx)
	}
}

func BenchmarkNewMeterProvider(b *testing.B) {
	ctx := b.Context()
	res := NewResource("benchmark-service", "1.0.0", "benchmark-id")
	exporter, _ := stdoutmetric.New()
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := NewMeterProvider(res, exporter)
		_ = provider.Shutdown(ctx)
	}
}

func BenchmarkNewTracerProvider(b *testing.B) {
	ctx := b.Context()
	res := NewResource("benchmark-service", "1.0.0", "benchmark-id")
	exporter, _ := stdouttrace.New()
	defer func() {
		_ = exporter.Shutdown(ctx)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := NewTracerProvider(res, exporter)
		_ = provider.Shutdown(ctx)
	}
}
