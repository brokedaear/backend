#!/bin/bash

# SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
#
# SPDX-License-Identifier: Apache-2.0

# Test script for OpenTelemetry gRPC integration

set -e

echo "🚀 Testing OpenTelemetry gRPC Integration"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if service is healthy
check_service() {
  local service_name=$1
  local url=$2
  local expected_code=${3:-200}

  echo -n "Checking $service_name... "

  if curl -s -o /dev/null -w "%{http_code}" "$url" | grep -q "$expected_code"; then
    echo -e "${GREEN}✓${NC}"
    return 0
  else
    echo -e "${RED}✗${NC}"
    return 1
  fi
}

# Function to wait for service
wait_for_service() {
  local service_name=$1
  local url=$2
  local max_attempts=30
  local attempt=1

  echo "Waiting for $service_name to be ready..."

  while [ $attempt -le $max_attempts ]; do
    if curl -s -o /dev/null "$url" 2>/dev/null; then
      echo -e "${GREEN}$service_name is ready!${NC}"
      return 0
    fi

    echo -n "."
    sleep 2
    attempt=$((attempt + 1))
  done

  echo -e "\n${RED}$service_name failed to start within timeout${NC}"
  return 1
}

echo "📋 Step 1: Starting telemetry stack..."
docker compose -f docker-compose.telemetry.yml up -d

echo "⏳ Step 2: Waiting for services to be ready..."
wait_for_service "Prometheus" "http://localhost:9090/-/ready"
wait_for_service "Grafana" "http://localhost:3000/api/health"
wait_for_service "Jaeger" "http://localhost:16686/"
wait_for_service "OpenTelemetry Collector" "http://localhost:8888/metrics"

echo "🔍 Step 3: Checking service health..."
check_service "Prometheus" "http://localhost:9090/-/ready"
check_service "Grafana" "http://localhost:3000/api/health"
check_service "Jaeger" "http://localhost:16686/"
check_service "OpenTelemetry Collector" "http://localhost:8888/metrics"
check_service "OpenTelemetry Collector gRPC" "http://localhost:4317/" 000

echo "🏗️  Step 4: Building example application..."
if [ ! -f "examples/telemetry_grpc_example.go" ]; then
  echo -e "${RED}Example application not found at examples/telemetry_grpc_example.go${NC}"
  exit 1
fi

go build -o telemetry-example ./examples/telemetry_grpc_example.go
echo -e "${GREEN}Application built successfully${NC}"

echo "🚀 Step 5: Starting example application..."
./telemetry-example &
APP_PID=$!

# Function to cleanup on exit
cleanup() {
  echo "🧹 Cleaning up..."
  kill $APP_PID 2>/dev/null || true
  docker compose -f docker-compose.telemetry.yml logs --tail=50
}
trap cleanup EXIT

echo "⏳ Step 6: Waiting for application to start..."
sleep 5

# Check if app is running
if ! kill -0 $APP_PID 2>/dev/null; then
  echo -e "${RED}Application failed to start${NC}"
  exit 1
fi

echo -e "${GREEN}Application is running (PID: $APP_PID)${NC}"

echo "📊 Step 7: Testing telemetry data flow..."

# Wait a bit for telemetry data to flow
echo "Waiting for telemetry data to be generated..."
sleep 30

# Check if metrics are available in Prometheus
echo "Checking for gRPC health metrics in Prometheus..."
if curl -s "http://localhost:9090/api/v1/query?query=grpc_health_status" | jq -r '.data.result[] | select(.metric.service_name) | .value[1]' | grep -q "1"; then
  echo -e "${GREEN}✓ gRPC health metrics found${NC}"
else
  echo -e "${YELLOW}⚠ gRPC health metrics not found (may take time to appear)${NC}"
fi

# Check for custom metrics
echo "Checking for custom application metrics..."
if curl -s "http://localhost:9090/api/v1/query?query=database_connection_health" | jq -r '.data.result' | grep -q "database_connection_health"; then
  echo -e "${GREEN}✓ Custom application metrics found${NC}"
else
  echo -e "${YELLOW}⚠ Custom application metrics not found${NC}"
fi

# Check if traces are available in Jaeger
echo "Checking for traces in Jaeger..."
if curl -s "http://localhost:16686/api/services" | jq -r '.data[]' | grep -q "brokedaear-backend"; then
  echo -e "${GREEN}✓ Traces found in Jaeger${NC}"
else
  echo -e "${YELLOW}⚠ Traces not found in Jaeger${NC}"
fi

echo "🎯 Step 8: Testing gRPC health endpoints..."

# Install grpc_health_probe if not available
if ! command -v grpc_health_probe &>/dev/null; then
  echo "Installing grpc_health_probe..."
  GRPC_HEALTH_PROBE_VERSION=v0.4.15
  wget -qO/tmp/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64
  chmod +x /tmp/grpc_health_probe
  GRPC_HEALTH_PROBE_CMD="/tmp/grpc_health_probe"
else
  GRPC_HEALTH_PROBE_CMD="grpc_health_probe"
fi

# Test overall server health
echo "Testing overall server health..."
if $GRPC_HEALTH_PROBE_CMD -addr=localhost:8080; then
  echo -e "${GREEN}✓ Overall server health check passed${NC}"
else
  echo -e "${RED}✗ Overall server health check failed${NC}"
fi

# Test specific service health
echo "Testing database service health..."
if $GRPC_HEALTH_PROBE_CMD -addr=localhost:8080 -service=database-service; then
  echo -e "${GREEN}✓ Database service health check passed${NC}"
else
  echo -e "${YELLOW}⚠ Database service health check failed (may be NOT_SERVING)${NC}"
fi

echo "📈 Step 9: Generating sample queries..."

# Generate some Prometheus queries
echo "Sample Prometheus queries you can run:"
echo "  - Service health: http://localhost:9090/graph?g0.expr=grpc_health_status"
echo "  - Health check rate: http://localhost:9090/graph?g0.expr=rate(grpc_health_checks_total[5m])"
echo "  - Response time p95: http://localhost:9090/graph?g0.expr=histogram_quantile(0.95,%20rate(sample_response_time_bucket[5m]))"

echo "📊 Sample Grafana URLs:"
echo "  - Main dashboard: http://localhost:3000/d/1/bde-grpc-health-performance"
echo "  - Login: admin/admin"

echo "🔍 Sample Jaeger URLs:"
echo "  - Service traces: http://localhost:16686/search?service=brokedaear-backend"

echo -e "\n${GREEN}🎉 Testing completed!${NC}"
echo -e "\n📋 Summary:"
echo "  - Prometheus: http://localhost:9090"
echo "  - Grafana: http://localhost:3000 (admin/admin)"
echo "  - Jaeger: http://localhost:16686"
echo "  - Your app: localhost:8080 (gRPC)"
echo ""
echo "The application will continue running. Press Ctrl+C to stop."

# Keep the script running
wait $APP_PID

