// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"

	"backend.brokedaear.com"
	"backend.brokedaear.com/internal/common/telemetry"
	"backend.brokedaear.com/internal/common/utils/loggers"
	"backend.brokedaear.com/internal/core/server"
)

func main() {
	ctx := context.Background()

	// Setup telemetry with gRPC exporter to local OpenTelemetry Collector
	telConfig := &telemetry.Config{
		ServiceName:    "brokedabackend",
		ServiceVersion: "0.1.0",
		ServiceID:      "local-instance-1",
		ExporterConfig: *telemetry.NewExporterConfig(
			telemetry.WithType(telemetry.ExporterTypeGRPC),
			telemetry.WithEndpoint("localhost:4317"), // OpenTelemetry Collector gRPC endpoint
		),
	}

	tel, err := telemetry.New(ctx, telConfig)
	if err != nil {
		log.Fatal("Failed to setup telemetry:", err)
	}
	defer tel.Close()

	// Setup logger with telemetry integration
	logger, err := loggers.NewZap(&loggers.ZapConfig{
		Env:                backend.EnvDevelopment,
		OtelServiceName:    "brokedaear-backend",
		OtelLoggerProvider: tel.LoggerProvider(), // Integrate with OpenTelemetry logger
		CustomZapper:       nil,
		WithTelemetry:      true, // Enable telemetry integration
	})
	if err != nil {
		fmt.Printf("failed to setup logger %v", err)
		return
	}

	// Setup gRPC server with telemetry enabled
	serverConfig := &server.Config{
		Addr:      server.Address("0.0.0.0"),
		Port:      server.Port(8080),
		Env:       backend.EnvDevelopment,
		Version:   server.Version("1.0.0"),
		Telemetry: true, // Enables automatic gRPC tracing
	}

	srv, err := server.NewGRPCServer(ctx, logger, serverConfig)
	if err != nil {
		fmt.Printf("failed to setup grpc server %v", err)
		return
	}
	defer srv.Close()

	// Setup health monitoring with telemetry
	setupHealthMonitoring(srv, tel, logger)

	// Start server
	go func() {
		logger.Info("Starting gRPC server", "port", 8080)
		if err := srv.ListenAndServe(ctx); err != nil {
			logger.Error("Server error", "error", err)
		}
	}()

	// Generate some sample telemetry data
	go generateSampleTelemetry(tel, srv, logger)

	// Wait for shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info("Shutting down...")
}

func setupHealthMonitoring(srv server.GRPCServer, tel telemetry.Telemetry, logger server.Logger) {
	// Create custom metrics for monitoring
	dbHealthGauge, err := tel.Gauge(telemetry.Metric{
		Name:        "database_connection_health",
		Unit:        "{status}",
		Description: "Database connection health status (1=healthy, 0=unhealthy)",
	})
	if err != nil {
		logger.Error("Failed to create database health gauge", "error", err)
		return
	}

	cacheHealthGauge, err := tel.Gauge(telemetry.Metric{
		Name:        "cache_connection_health",
		Unit:        "{status}",
		Description: "Cache connection health status (1=healthy, 0=unhealthy)",
	})
	if err != nil {
		logger.Error("Failed to create cache health gauge", "error", err)
		return
	}

	const healthMonitorInterval = 30 * time.Second

	go func() {
		for {
			<-time.After(healthMonitorInterval)
			ctx := context.Background()

			// Start a trace for health monitoring
			ctx, span := tel.TraceStart(ctx, "health_monitoring_cycle")

			// Simulate database health check
			dbHealthy := simulateDatabaseHealthCheck(ctx, tel)
			dbStatus := grpc_health_v1.HealthCheckResponse_NOT_SERVING
			if dbHealthy {
				dbStatus = grpc_health_v1.HealthCheckResponse_SERVING
			}

			// Record database health metric
			var dbHealthValue int64
			if dbHealthy {
				dbHealthValue = 1
			}
			dbHealthGauge.Record(ctx, dbHealthValue)

			// Simulate cache health check
			cacheHealthy := simulateCacheHealthCheck(ctx, tel)
			cacheStatus := grpc_health_v1.HealthCheckResponse_NOT_SERVING
			if cacheHealthy {
				cacheStatus = grpc_health_v1.HealthCheckResponse_SERVING
			}

			// Record cache health metric
			cacheHealthValue := int64(0)
			if cacheHealthy {
				cacheHealthValue = 1
			}
			cacheHealthGauge.Record(ctx, cacheHealthValue)

			// Update service health statuses
			srv.SetHealthStatus("database-service", dbStatus)
			srv.SetHealthStatus("cache-service", cacheStatus)

			// Overall health depends on all dependencies
			overallStatus := grpc_health_v1.HealthCheckResponse_SERVING
			if !dbHealthy || !cacheHealthy {
				overallStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
			}
			srv.SetHealthStatus("", overallStatus)

			logger.Info("Health check completed",
				"database_healthy", dbHealthy,
				"cache_healthy", cacheHealthy,
				"overall_status", overallStatus.String())

			span.End()
		}
	}()
}

func simulateDatabaseHealthCheck(ctx context.Context, tel telemetry.Telemetry) bool {
	_, span := tel.TraceStart(ctx, "database_health_check")
	defer span.End()

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// 90% chance of being healthy
	healthy := time.Now().Unix()%10 != 0
	return healthy
}

func simulateCacheHealthCheck(ctx context.Context, tel telemetry.Telemetry) bool {
	_, span := tel.TraceStart(ctx, "cache_health_check")
	defer span.End()

	// Simulate some work
	time.Sleep(5 * time.Millisecond)

	// 95% chance of being healthy
	healthy := time.Now().Unix()%20 != 0
	return healthy
}

func generateSampleTelemetry(tel telemetry.Telemetry, _ server.GRPCServer, logger server.Logger) {
	// Create sample metrics
	requestCounter, err := tel.UpDownCounter(telemetry.Metric{
		Name:        "sample_requests_total",
		Unit:        "{count}",
		Description: "Total number of sample requests",
	})
	if err != nil {
		logger.Error("Failed to create request counter", "error", err)
		return
	}

	responseTimeHistogram, err := tel.Histogram(telemetry.Metric{
		Name:        "sample_response_time",
		Unit:        "ms",
		Description: "Sample response time in milliseconds",
	})
	if err != nil {
		logger.Error("Failed to create response time histogram", "error", err)
		return
	}

	workDuration := time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond

	for {
		<-time.After(5 * time.Second)
		ctx := context.Background()

		// Generate sample traces and metrics
		ctx, span := tel.TraceStart(ctx, "sample_operation")

		// Simulate some work
		select {
		case <-time.After(workDuration):
			// Record metrics
			requestCounter.Add(ctx, 1)
			responseTimeHistogram.Record(ctx, workDuration.Milliseconds())

			span.End()

			logger.Debug("Generated sample telemetry data",
				"duration_ms", workDuration.Milliseconds())
		default:
		}
	}
}
