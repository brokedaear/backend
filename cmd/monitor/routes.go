// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func routePrometheusMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
