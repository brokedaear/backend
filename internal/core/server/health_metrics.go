// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc/health/grpc_health_v1"

	"backend.brokedaear.com/internal/common/telemetry"
)

// HealthMetrics provides telemetry integration for health checks.
type HealthMetrics struct {
	healthStatusGauge    otelmetric.Int64Gauge
	healthCheckCounter   otelmetric.Int64UpDownCounter
	healthWatcherGauge   otelmetric.Int64UpDownCounter
	telemetry            telemetry.Telemetry
}

// NewHealthMetrics creates a new health metrics instance.
func NewHealthMetrics(tel telemetry.Telemetry) (*HealthMetrics, error) {
	healthStatusGauge, err := tel.Gauge(telemetry.MetricGRPCHealthStatus)
	if err != nil {
		return nil, err
	}

	healthCheckCounter, err := tel.UpDownCounter(telemetry.MetricGRPCHealthChecksTotal)
	if err != nil {
		return nil, err
	}

	healthWatcherGauge, err := tel.UpDownCounter(telemetry.MetricGRPCHealthWatchers)
	if err != nil {
		return nil, err
	}

	return &HealthMetrics{
		healthStatusGauge:    healthStatusGauge,
		healthCheckCounter:   healthCheckCounter,
		healthWatcherGauge:   healthWatcherGauge,
		telemetry:            tel,
	}, nil
}

// RecordHealthStatus records the current health status for a service.
func (hm *HealthMetrics) RecordHealthStatus(ctx context.Context, service string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	attributes := otelmetric.WithAttributes(
		attribute.String("service", service),
		attribute.String("status", status.String()),
	)
	
	hm.healthStatusGauge.Record(ctx, int64(status), attributes)
}

// RecordHealthCheck records a health check request.
func (hm *HealthMetrics) RecordHealthCheck(ctx context.Context, service string, status grpc_health_v1.HealthCheckResponse_ServingStatus, success bool) {
	attributes := otelmetric.WithAttributes(
		attribute.String("service", service),
		attribute.String("status", status.String()),
		attribute.Bool("success", success),
	)
	
	hm.healthCheckCounter.Add(ctx, 1, attributes)
}

// RecordWatcherChange records a change in the number of health watchers.
func (hm *HealthMetrics) RecordWatcherChange(ctx context.Context, service string, delta int64) {
	attributes := otelmetric.WithAttributes(
		attribute.String("service", service),
	)
	
	hm.healthWatcherGauge.Add(ctx, delta, attributes)
}