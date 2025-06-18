# OpenTelemetry gRPC Integration Setup

This guide shows how to set up OpenTelemetry with gRPC exporter, connecting to Grafana/Prometheus on localhost.

## Quick Start

### 1. Start the Telemetry Stack

```bash
# Start OpenTelemetry Collector, Prometheus, Grafana, and Jaeger
docker-compose -f docker-compose.telemetry.yml up -d

# Check all services are running
docker-compose -f docker-compose.telemetry.yml ps
```

### 2. Build and Run Your Application

```bash
# Build the example application
go build -o telemetry-example ./examples/telemetry_grpc_example/

# Run with gRPC telemetry enabled
./telemetry-example
```

### 3. Access the Dashboards

- **Grafana**: <http://localhost:3000> (admin/admin)
- **Prometheus**: <http://localhost:9090>
- **Jaeger**: <http://localhost:16686>
- **OpenTelemetry Collector Metrics**: <http://localhost:8888/metrics>

## Architecture Overview

```
Your App (port 8080)
    ↓ gRPC telemetry
OpenTelemetry Collector (port 4317)
    ↓ metrics          ↓ traces
Prometheus (port 9090) → Grafana (port 3000)
                       ↓
                   Jaeger (port 16686)
```

## Configuration Details

### Application Configuration

Your app sends telemetry data via gRPC to the OpenTelemetry Collector:

```go
telConfig := &telemetry.Config{
    ServiceName:    "brokedaear-backend",
    ServiceVersion: "1.0.0",
    ServiceID:      "local-instance-1",
    ExporterConfig: telemetry.NewExporterConfig(
        telemetry.ExporterTypeGRPC,
        "http://localhost:4317", // Collector gRPC endpoint
        true,                    // Insecure for local
        nil,                     // No auth headers
    ),
}
```

### OpenTelemetry Collector

The collector receives gRPC telemetry data and exports:

- **Metrics** → Prometheus (port 8889)
- **Traces** → Jaeger (port 14250)

### Prometheus Scraping

Prometheus scrapes metrics from:

- OpenTelemetry Collector (port 8889)
- Your app directly (port 8080) - if you expose /metrics endpoint

## Available Metrics

### gRPC Health Metrics

- `grpc_health_status` - Current health status of services
- `grpc_health_checks_total` - Total health check requests
- `grpc_health_watchers` - Active health check watchers

### Custom Application Metrics

- `database_connection_health` - Database connection status
- `cache_connection_health` - Cache connection status
- `sample_requests_total` - Sample request counter
- `sample_response_time` - Sample response time histogram

### Automatic gRPC Metrics (from otelgrpc)

- `grpc_server_handling_seconds` - Request duration
- `grpc_server_requests_total` - Request count
- `grpc_server_responses_total` - Response count

## Testing the Setup

### 1. Generate Sample Data

The example application automatically generates:

- Health check data every 30 seconds
- Sample telemetry data every 5 seconds
- gRPC traces for all operations

### 2. Verify Data Flow

```bash
# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check available metrics
curl http://localhost:9090/api/v1/label/__name__/values

# Query health status
curl "http://localhost:9090/api/v1/query?query=grpc_health_status"
```

### 3. View in Grafana

1. Go to <http://localhost:3000>
2. Login with admin/admin
3. Navigate to "BdE gRPC Health & Performance" dashboard
4. You should see:
   - Service health status panels
   - Health check request rates
   - Response time distributions
   - Database/cache health indicators

### 4. View Traces in Jaeger

1. Go to <http://localhost:16686>
2. Select "brokedaear-backend" service
3. Click "Find Traces"
4. Explore distributed traces showing:
   - Health monitoring cycles
   - Database/cache health checks
   - Sample operations

## Troubleshooting

### Common Issues

1. **No data in Grafana**

   ```bash
   # Check collector logs
   docker-compose -f docker-compose.telemetry.yml logs otel-collector

   # Check if metrics are being exported
   curl http://localhost:8889/metrics
   ```

2. **Application can't connect to collector**

   ```bash
   # Check collector is listening
   netstat -ln | grep 4317

   # Test gRPC connection
   grpcurl -plaintext localhost:4317 list
   ```

3. **Prometheus not scraping**

   ```bash
   # Check Prometheus targets
   curl http://localhost:9090/api/v1/targets

   # Check Prometheus logs
   docker-compose -f docker-compose.telemetry.yml logs prometheus
   ```

### Health Check Commands

```bash
# Test gRPC health endpoint directly
grpc_health_probe -addr=localhost:8080

# Check specific service health
grpc_health_probe -addr=localhost:8080 -service=database-service

# Watch health status changes
grpc_health_probe -addr=localhost:8080 -service="" -watch
```

## Environment Variables

You can override configuration with environment variables:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4317"
export OTEL_EXPORTER_OTLP_INSECURE="true"
export OTEL_SERVICE_NAME="brokedaear-backend"
export OTEL_SERVICE_VERSION="1.0.0"
```

## Production Considerations

For production deployments:

1. **Enable TLS** for the gRPC exporter
2. **Add authentication** headers
3. **Configure proper retention** in Prometheus
4. **Set up alerting rules** based on health metrics
5. **Use persistent volumes** for data storage
6. **Configure resource limits** for all containers

## Cleanup

```bash
# Stop all services
docker-compose -f docker-compose.telemetry.yml down

# Remove volumes (careful - this deletes all data)
docker-compose -f docker-compose.telemetry.yml down -v
```

