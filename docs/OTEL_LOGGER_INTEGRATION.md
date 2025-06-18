# OpenTelemetry Logger Provider Integration

This guide shows how to integrate your application logs with OpenTelemetry using the `OtelLoggerProvider` for complete observability (logs, metrics, traces).

## How It Works

The integration works by:

1. **Creating a telemetry instance** with log, metric, and trace providers
2. **Passing the log provider** to your Zap logger configuration  
3. **Automatically correlating logs** with traces and spans
4. **Exporting logs** to your observability backend alongside metrics/traces

## Architecture

```
Your App
    ↓ Structured Logs (with trace correlation)
Zap Logger (with OpenTelemetry bridge)
    ↓ OpenTelemetry Log Records
OpenTelemetry Collector
    ↓ Logs, Metrics, Traces
Observability Backend (Grafana/Loki/Jaeger/Prometheus)
```

## Implementation

### 1. Basic Setup

```go
package main

import (
    "context"
    
    "backend.brokedaear.com"
    "backend.brokedaear.com/internal/common/telemetry"
    "backend.brokedaear.com/internal/common/utils/loggers"
)

func main() {
    ctx := context.Background()

    // 1. Setup telemetry first
    telConfig := &telemetry.Config{
        ServiceName:    "my-service",
        ServiceVersion: "1.0.0",
        ServiceID:      "instance-1",
        ExporterConfig: telemetry.NewExporterConfig(
            telemetry.ExporterTypeGRPC,
            "http://localhost:4317",
            true,
            nil,
        ),
    }

    tel, err := telemetry.New(ctx, telConfig)
    if err != nil {
        panic(err)
    }
    defer tel.Close()

    // 2. Setup logger with OpenTelemetry integration
    logger, err := loggers.NewZap(&loggers.ZapConfig{
        Env:                backend.EnvDevelopment,
        OtelServiceName:    "my-service",
        OtelLoggerProvider: tel.LoggerProvider(), // ← This is the key integration
        CustomZapper:       nil,
        WithTelemetry:      true,
    })
    if err != nil {
        panic(err)
    }

    // 3. Use structured logging with automatic trace correlation
    useStructuredLogging(ctx, tel, logger)
}
```

### 2. Structured Logging with Trace Correlation

```go
func useStructuredLogging(ctx context.Context, tel telemetry.Telemetry, logger interface{}) {
    // Cast to logger interface
    log := logger.(interface {
        Info(msg string, args ...any)
        Error(msg string, args ...any)
        // ... other methods
    })

    // Start a trace - logs within this span will be correlated
    ctx, span := tel.TraceStart(ctx, "business_operation")
    defer span.End()

    // These logs will automatically include trace/span IDs
    log.Info("Processing user request",
        "user_id", "12345",
        "operation", "update_profile",
        "request_id", "req-abc-123")

    // Simulate some work
    if err := processUserData(); err != nil {
        // Error logs include full context
        log.Error("Failed to process user data",
            "error", err.Error(),
            "user_id", "12345",
            "retry_count", 3)
        
        // Set span status for trace correlation
        span.RecordError(err)
        span.SetStatus(codes.Error, "Processing failed")
        return
    }

    log.Info("User profile updated successfully",
        "user_id", "12345",
        "fields_updated", []string{"email", "preferences"})
}
```

### 3. Advanced Integration with Health Monitoring

```go
func setupHealthMonitoringWithLogs(srv server.GRPCServer, tel telemetry.Telemetry, logger interface{}) {
    log := logger.(interface {
        Info(msg string, args ...any)
        Warn(msg string, args ...any)
        Error(msg string, args ...any)
    })

    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                ctx := context.Background()
                
                // Start health check trace
                ctx, span := tel.TraceStart(ctx, "health_check_cycle")
                
                // Check database health with structured logging
                dbHealthy := checkDatabaseHealth(ctx, tel, log)
                
                // Log health status with structured data
                log.Info("Health check completed",
                    "component", "database",
                    "healthy", dbHealthy,
                    "check_duration_ms", 25,
                    "timestamp", time.Now().Unix())

                // Update gRPC health status
                if dbHealthy {
                    srv.SetHealthStatus("database", grpc_health_v1.HealthCheckResponse_SERVING)
                } else {
                    srv.SetHealthStatus("database", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
                    log.Warn("Database health check failed",
                        "component", "database",
                        "action", "marking_service_unhealthy")
                }

                span.End()
            }
        }
    }()
}

func checkDatabaseHealth(ctx context.Context, tel telemetry.Telemetry, log interface{}) bool {
    // Create child span for database check
    ctx, span := tel.TraceStart(ctx, "database_health_check")
    defer span.End()

    // Log the health check attempt
    logger := log.(interface {
        Debug(msg string, args ...any)
        Error(msg string, args ...any)
    })

    logger.Debug("Starting database health check",
        "timeout_ms", 5000,
        "connection_pool", "primary")

    // Simulate health check
    healthy := performHealthCheck()
    
    if !healthy {
        logger.Error("Database health check failed",
            "error", "connection timeout",
            "endpoint", "postgres://localhost:5432")
        span.SetStatus(codes.Error, "Health check failed")
    }

    return healthy
}
```

### 4. Business Event Logging

```go
func logBusinessEvents(ctx context.Context, tel telemetry.Telemetry, logger interface{}) {
    log := logger.(interface {
        Info(msg string, args ...any)
    })

    // Business events with rich context
    ctx, span := tel.TraceStart(ctx, "user_signup")
    defer span.End()

    // Log structured business event
    log.Info("User signup initiated",
        "event_type", "user_signup_started",
        "user_email", "user@example.com",
        "signup_method", "email",
        "referral_source", "organic",
        "timestamp", time.Now().Unix())

    // ... business logic ...

    log.Info("User signup completed",
        "event_type", "user_signup_completed", 
        "user_id", "12345",
        "user_email", "user@example.com",
        "activation_required", true,
        "welcome_email_sent", true)
}
```

## What You Get

### 1. **Automatic Trace Correlation**
- Every log within a trace span automatically includes trace/span IDs
- You can jump from logs to traces in Grafana/Jaeger
- Full request flow visibility across services

### 2. **Structured Log Export**
Logs are exported in OpenTelemetry format with:
```json
{
  "timestamp": "2025-01-18T10:30:00Z",
  "level": "INFO",
  "message": "User profile updated successfully",
  "service.name": "my-service",
  "service.version": "1.0.0",
  "trace_id": "abc123...",
  "span_id": "def456...",
  "attributes": {
    "user_id": "12345",
    "fields_updated": ["email", "preferences"]
  }
}
```

### 3. **Multiple Export Destinations**
Configure where logs go in the OpenTelemetry Collector:
- **File export**: `/tmp/otel-logs.json`
- **Loki**: For log aggregation and querying
- **Elasticsearch**: For full-text search
- **Cloud providers**: AWS CloudWatch, GCP Cloud Logging, etc.

### 4. **Grafana Integration**
- **Logs panel**: View logs alongside metrics and traces
- **Trace to logs**: Click on trace spans to see related logs
- **Log correlation**: Filter logs by trace ID, service, user ID, etc.

## Configuration Examples

### Development Setup
```go
logger, err := loggers.NewZap(&loggers.ZapConfig{
    Env:                backend.EnvDevelopment,
    OtelServiceName:    "my-service",
    OtelLoggerProvider: tel.LoggerProvider(),
    WithTelemetry:      true, // Enable for development
})
```

### Production Setup
```go
logger, err := loggers.NewZap(&loggers.ZapConfig{
    Env:                backend.EnvProduction,
    OtelServiceName:    "my-service", 
    OtelLoggerProvider: tel.LoggerProvider(),
    WithTelemetry:      true, // Enable for production observability
})
```

### Custom Writer Setup
```go
customWriter := loggers.NewCustomZapWriter(
    "/app/logs/app.log",
    "file",
    "caller",
    &fileWriter{},
)

logger, err := loggers.NewZap(&loggers.ZapConfig{
    Env:                backend.EnvProduction,
    OtelServiceName:    "my-service",
    OtelLoggerProvider: tel.LoggerProvider(),
    CustomZapper:       customWriter, // Local file + OpenTelemetry
    WithTelemetry:      true,
})
```

## Collector Configuration

Update your OpenTelemetry Collector to handle logs:

```yaml
# configs/otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

exporters:
  # Export logs to file
  file:
    path: /tmp/otel-logs.json
    
  # Export logs to Loki
  loki:
    endpoint: http://loki:3100/loki/api/v1/push
    
service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [file, loki] # Choose your exporters
```

## Querying and Alerting

### Grafana Log Queries
```logql
# All logs from your service
{service_name="my-service"}

# Error logs only  
{service_name="my-service"} |= "ERROR"

# Logs for specific user
{service_name="my-service"} | json | user_id="12345"

# Logs within a trace
{service_name="my-service"} | json | trace_id="abc123..."
```

### Prometheus Alerts on Logs
You can create metrics from logs and alert on them:
```yaml
# Alert on error rate
- alert: HighErrorRate
  expr: rate(log_entries{level="ERROR"}[5m]) > 0.1
  annotations:
    summary: "High error rate detected in {{ $labels.service_name }}"
```

## Performance Considerations

1. **Log Level**: Use appropriate log levels (DEBUG only in development)
2. **Sampling**: Consider log sampling for high-throughput services
3. **Buffering**: OpenTelemetry batches logs for efficiency
4. **Resource Limits**: Set memory limits in the collector configuration

This integration provides complete observability where logs, metrics, and traces work together to give you full visibility into your application's behavior.