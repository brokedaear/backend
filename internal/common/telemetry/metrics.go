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
