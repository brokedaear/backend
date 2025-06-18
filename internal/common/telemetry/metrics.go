// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

// Metric represents a metric that can be collected by the server.
type Metric struct {
	Name        string
	Unit        string
	Description string
}

// MetricRequestDurationMillis is a metric that measures the latency of HTTP
// requests processed by the server, in milliseconds.
var MetricRequestDurationMillis = Metric{ //nolint:gochecknoglobals // makes more sense like this.
	Name:        "request_duration_millis",
	Unit:        "ms",
	Description: "Measures the latency of HTTP requests processed by the server, in milliseconds.",
}

// MetricRequestsInFlight is a metric that measures the number of requests
// currently being processed by the server.
var MetricRequestsInFlight = Metric{ //nolint:gochecknoglobals // makes more sense like this.
	Name:        "requests_inflight",
	Unit:        "{count}",
	Description: "Measures the number of requests currently being processed by the server.",
}

// MetricGRPCHealthStatus is a metric that tracks the current health status of gRPC services.
var MetricGRPCHealthStatus = Metric{ //nolint:gochecknoglobals // makes more sense like this.
	Name:        "grpc_health_status",
	Unit:        "{status}",
	Description: "Current health status of gRPC services (0=UNKNOWN, 1=SERVING, 2=NOT_SERVING, 3=SERVICE_UNKNOWN).",
}

// MetricGRPCHealthChecksTotal is a metric that counts health check requests.
var MetricGRPCHealthChecksTotal = Metric{ //nolint:gochecknoglobals // makes more sense like this.
	Name:        "grpc_health_checks_total",
	Unit:        "{count}",
	Description: "Total number of health check requests received.",
}

// MetricGRPCHealthWatchers is a metric that tracks active health check watchers.
var MetricGRPCHealthWatchers = Metric{ //nolint:gochecknoglobals // makes more sense like this.
	Name:        "grpc_health_watchers",
	Unit:        "{count}",
	Description: "Number of active health check watchers.",
}
