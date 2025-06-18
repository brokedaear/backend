// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry_test

import (
	"testing"

	"backend.brokedaear.com/internal/common/telemetry"
	"backend.brokedaear.com/pkg/assert"
)

func TestMetric(t *testing.T) {
	metric := telemetry.Metric{
		Name:        "test_metric",
		Unit:        "ms",
		Description: "A test metric",
	}

	assert.Equal(t, metric.Name, "test_metric")
	assert.Equal(t, metric.Unit, "ms")
	assert.Equal(t, metric.Description, "A test metric")
}

func TestMetricRequestDurationMillis(t *testing.T) {
	metric := telemetry.MetricRequestDurationMillis

	assert.Equal(t, metric.Name, "request_duration_millis")
	assert.Equal(t, metric.Unit, "ms")
	assert.Equal(
		t,
		metric.Description,
		"Measures the latency of HTTP requests processed by the server, in milliseconds.",
	)
}

func TestMetricRequestsInFlight(t *testing.T) {
	metric := telemetry.MetricRequestsInFlight

	assert.Equal(t, metric.Name, "requests_inflight")
	assert.Equal(t, metric.Unit, "{count}")
	assert.Equal(
		t,
		metric.Description,
		"Measures the number of requests currently being processed by the server.",
	)
}

func TestMetricStructFields(t *testing.T) {
	tests := []struct {
		name   string
		metric telemetry.Metric
	}{
		{
			name:   "request duration metric",
			metric: telemetry.MetricRequestDurationMillis,
		},
		{
			name:   "requests in flight metric",
			metric: telemetry.MetricRequestsInFlight,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Verify all fields are non-empty
				assert.NotEqual(t, tt.metric.Name, "")
				assert.NotEqual(t, tt.metric.Unit, "")
				assert.NotEqual(t, tt.metric.Description, "")

				// Verify Name field doesn't contain spaces (following metric naming conventions)
				for _, char := range tt.metric.Name {
					if char == ' ' {
						t.Errorf("metric name '%s' should not contain spaces", tt.metric.Name)
					}
				}

				// Verify Unit field follows proper format
				assert.True(t, len(tt.metric.Unit) > 0)

				// Verify Description is meaningful
				assert.True(t, len(tt.metric.Description) > len(tt.metric.Name))
			},
		)
	}
}

func TestMetricZeroValues(t *testing.T) {
	var metric telemetry.Metric

	assert.Equal(t, metric.Name, "")
	assert.Equal(t, metric.Unit, "")
	assert.Equal(t, metric.Description, "")
}

func TestMetricComparison(t *testing.T) {
	metric1 := telemetry.Metric{
		Name:        "test_metric",
		Unit:        "ms",
		Description: "A test metric",
	}

	metric2 := telemetry.Metric{
		Name:        "test_metric",
		Unit:        "ms",
		Description: "A test metric",
	}

	metric3 := telemetry.Metric{
		Name:        "different_metric",
		Unit:        "ms",
		Description: "A test metric",
	}

	// Test equality
	assert.Equal(t, metric1, metric2)

	// Test inequality
	assert.NotEqual(t, metric1, metric3)
}

func BenchmarkMetricAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = telemetry.MetricRequestDurationMillis.Name
		_ = telemetry.MetricRequestDurationMillis.Unit
		_ = telemetry.MetricRequestDurationMillis.Description
	}
}

func BenchmarkMetricCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = telemetry.Metric{
			Name:        "test_metric",
			Unit:        "ms",
			Description: "A test metric for benchmarking",
		}
	}
}
