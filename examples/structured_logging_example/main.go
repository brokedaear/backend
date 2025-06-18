// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend.brokedaear.com"
	"backend.brokedaear.com/internal/common/telemetry"
	"backend.brokedaear.com/internal/common/utils/loggers"
)

func main() {
	ctx := context.Background()

	// Setup telemetry with gRPC exporter
	telConfig := &telemetry.Config{
		ServiceName:    "logging-example",
		ServiceVersion: "1.0.0",
		ServiceID:      "local-instance-1",
		ExporterConfig: *telemetry.NewExporterConfig(
			telemetry.WithType(telemetry.ExporterTypeGRPC),
			telemetry.WithEndpoint("localhost:4317"),
		),
	}

	tel, err := telemetry.New(ctx, telConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to setup telemetry: %v", err))
	}
	defer tel.Close()

	// Setup logger with OpenTelemetry integration
	logger, err := loggers.NewZap(&loggers.ZapConfig{
		Env:                backend.EnvDevelopment,
		OtelServiceName:    "logging-example",
		OtelLoggerProvider: tel.LoggerProvider(), // This integrates logs with OpenTelemetry
		CustomZapper:       nil,
		WithTelemetry:      true,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to setup logger: %v", err))
	}

	// Demonstrate structured logging with telemetry correlation
	demonstrateStructuredLogging(ctx, tel, logger)
}

func demonstrateStructuredLogging(ctx context.Context, tel telemetry.Telemetry, logger any) {
	// Cast to the logger interface
	log := logger.(interface {
		Info(msg string, args ...any)
		Debug(msg string, args ...any)
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
		Sync() error
	})

	fmt.Println("🔍 Demonstrating structured logging with OpenTelemetry integration...")

	// Example 1: Basic structured logging
	log.Info("Application started",
		"service", "logging-example",
		"version", "1.0.0",
		"environment", "development")

	// Example 2: Logging within a trace span
	ctx, span := tel.TraceStart(ctx, "user_operation")

	// These logs will be correlated with the trace
	log.Info("Processing user request",
		"user_id", "12345",
		"operation", "user_profile_update",
		"request_id", "req-abc-123")

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	log.Debug("Database query executed",
		"query", "UPDATE users SET last_login = ? WHERE id = ?",
		"duration_ms", 25,
		"rows_affected", 1)

	// Example 3: Error logging with context
	simulateError := func() error {
		return errors.New("connection timeout")
	}

	if err := simulateError(); err != nil {
		log.Error("Database operation failed",
			"error", err.Error(),
			"retry_count", 3,
			"operation", "user_profile_update",
			"user_id", "12345")
	}

	span.End()

	// Example 4: Metrics and logging correlation
	counter, err := tel.UpDownCounter(telemetry.Metric{
		Name:        "user_operations_total",
		Unit:        "{count}",
		Description: "Total user operations processed",
	})
	if err == nil {
		counter.Add(ctx, 1)
		log.Info("Metrics recorded",
			"metric", "user_operations_total",
			"value", 1)
	}

	// Example 5: Different log levels
	log.Debug("Debug information for troubleshooting",
		"component", "user_service",
		"method", "updateProfile")

	log.Warn("Performance warning detected",
		"slow_query_threshold_ms", 100,
		"actual_duration_ms", 150,
		"query", "SELECT * FROM user_preferences")

	// Example 6: Structured logging for business events
	log.Info("Business event occurred",
		"event_type", "user_profile_updated",
		"user_id", "12345",
		"fields_updated", []string{"email", "preferences"},
		"timestamp", time.Now().Unix())

	fmt.Println("✅ Structured logging examples completed!")
	fmt.Println("📄 Check /tmp/otel-logs.json for exported logs")
	fmt.Println("🔍 Check Jaeger UI for trace correlation")

	// Ensure all logs are flushed
	if err := log.Sync(); err != nil {
		fmt.Printf("Warning: Failed to sync logs: %v\n", err)
	}
}
