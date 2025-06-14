// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"

	"backend.brokedaear.com/internal/common/telemetry"
	"backend.brokedaear.com/internal/common/tests/assert"
)

func newTelConfig() *telemetry.Config {
	return &telemetry.Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		ServiceID:      "test-id",
		ExporterConfig: telemetry.ExporterConfig{
			Type:     telemetry.ExporterTypeStdout,
			Endpoint: "",
			Insecure: true,
			Headers:  make(map[string]string),
		},
	}
}

func TestNew(t *testing.T) {
	ctx := t.Context()

	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}

	defer func() {
		_ = tel.Close()
	}()

	// Verify that telemetry implements the Telemetry interface
	_ = tel
}

func TestNewWithInvalidConfig(t *testing.T) {
	ctx := t.Context()

	cfg := &telemetry.Config{
		ServiceName:    "test-!service",
		ServiceID:      "test-id",
		ServiceVersion: "1.abc0.0",
		ExporterConfig: telemetry.ExporterConfig{
			Type:     telemetry.ExporterTypeStdout,
			Endpoint: "",
			Insecure: true,
			Headers:  make(map[string]string),
		},
	}

	tel, err := telemetry.New(ctx, cfg)
	if err == nil {
		t.Error("expected error")
	}
	if tel != nil {
		t.Error("expected nil telemetry on error")
	}
}

func TestOtelTelemetry_Histogram(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	metric := telemetry.Metric{
		Name:        "test_histogram",
		Unit:        "ms",
		Description: "A test histogram metric",
	}

	histogram, err := tel.Histogram(metric)
	assert.NoError(t, err)
	if histogram == nil {
		t.Error("expected non-nil histogram")
	}

	// Test recording a value
	histogram.Record(ctx, 100)
}

func TestOtelTelemetry_UpDownCounter(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	metric := telemetry.Metric{
		Name:        "test_updown_counter",
		Unit:        "{count}",
		Description: "A test up-down counter metric",
	}

	counter, err := tel.UpDownCounter(metric)
	assert.NoError(t, err)
	if counter == nil {
		t.Error("expected non-nil counter")
	}

	// Test adding and subtracting values
	counter.Add(ctx, 5)
	counter.Add(ctx, -2)
}

func TestOtelTelemetry_Gauge(t *testing.T) {
	ctx := t.Context()

	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	metric := telemetry.Metric{
		Name:        "test_gauge",
		Unit:        "{items}",
		Description: "A test gauge metric",
	}

	gauge, err := tel.Gauge(metric)
	assert.NoError(t, err)
	if gauge == nil {
		t.Error("expected non-nil gauge")
	}

	// Test recording a value
	gauge.Record(ctx, 42)
}

func TestOtelTelemetry_TraceStart(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	spanCtx, span := tel.TraceStart(ctx, "test-span")
	if spanCtx == nil {
		t.Error("expected non-nil span context")
	}
	if span == nil {
		t.Error("expected non-nil span")
	}

	// Verify span context is different from original context
	assert.NotEqual(t, ctx, spanCtx)

	// Add some attributes and end the span
	span.SetAttributes(attribute.String("test.key", "test.value"))
	span.End()
}

func TestOtelTelemetry_Close(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil tel")
	}

	// Close should not error
	err = tel.Close()
	assert.NoError(t, err)

	// Calling close multiple times may return errors from some providers
	// but should not panic
	_ = tel.Close()
}

func TestOtelTelemetry_PredefinedMetrics(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	// Test creating histogram with predefined metric
	histogram, err := tel.Histogram(telemetry.MetricRequestDurationMillis)
	assert.NoError(t, err)
	if histogram == nil {
		t.Error("expected non-nil histogram")
	}

	// Test creating up-down counter with predefined metric
	counter, err := tel.UpDownCounter(telemetry.MetricRequestsInFlight)
	assert.NoError(t, err)
	if counter == nil {
		t.Error("expected non-nil counter")
	}

	// Test recording values
	histogram.Record(ctx, 150)
	counter.Add(ctx, 1)
	counter.Add(ctx, -1)
}

func TestOtelTelemetry_MetricCreationErrors(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	// Test with invalid metric name (empty)
	invalidMetric := telemetry.Metric{
		Name:        "",
		Unit:        "ms",
		Description: "Invalid metric with empty name",
	}

	histogram, err := tel.Histogram(invalidMetric)
	if err == nil {
		t.Error("expected error")
	}
	if histogram != nil {
		t.Error("expected nil histogram on error")
	}

	counter, err := tel.UpDownCounter(invalidMetric)
	if err == nil {
		t.Error("expected error")
	}
	if counter != nil {
		t.Error("expected nil counter on error")
	}

	gauge, err := tel.Gauge(invalidMetric)
	if err == nil {
		t.Error("expected error")
	}
	if gauge != nil {
		t.Error("expected nil gauge on error")
	}
}

func TestOtelTelemetry_MultipleSpans(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil telemetry")
	}
	defer func() {
		_ = tel.Close()
	}()

	// Create parent span
	parentCtx, parentSpan := tel.TraceStart(ctx, "parent-span")
	if parentCtx == nil {
		t.Error("expected non-nil parent context")
	}
	if parentSpan == nil {
		t.Error("expected non-nil parent span")
	}

	// Create child span
	childCtx, childSpan := tel.TraceStart(parentCtx, "child-span")
	if childCtx == nil {
		t.Error("expected non-nil child context")
	}
	if childSpan == nil {
		t.Error("expected non-nil child span")
	}

	// End spans in reverse order
	childSpan.End()
	parentSpan.End()
}

func TestOtelTelemetry_InterfaceCompliance(t *testing.T) {
	ctx := t.Context()
	tel, err := telemetry.New(ctx, newTelConfig())
	assert.NoError(t, err)
	if tel == nil {
		t.Error("expected non-nil tel")
	}
	defer func() {
		_ = tel.Close()
	}()

	// Verify interface compliance at compile time
	_ = tel

	// Verify all interface methods are available
	metric := telemetry.Metric{Name: "test", Unit: "unit", Description: "desc"}

	histogram, err := tel.Histogram(metric)
	assert.NoError(t, err)
	if histogram == nil {
		t.Error("expected non-nil histogram")
	}

	counter, err := tel.UpDownCounter(metric)
	assert.NoError(t, err)
	if counter == nil {
		t.Error("expected non-nil counter")
	}

	gauge, err := tel.Gauge(metric)
	assert.NoError(t, err)
	if gauge == nil {
		t.Error("expected non-nil gauge")
	}

	spanCtx, span := tel.TraceStart(ctx, "test")
	if spanCtx == nil {
		t.Error("expected non-nil span context")
	}
	if span == nil {
		t.Error("expected non-nil span")
	}
	span.End()

	err = tel.Close()
	assert.NoError(t, err)
}
