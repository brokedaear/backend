// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package main

import "net/http"

// healthCheckHandler handles healthcheck requests. It uses writeJSON()
// to write system health data to the http.ResponseWriter stream.
func (a *app) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("received %s from %s", r.Method, r.RemoteAddr)

	res := jsonWrap{
		"status": "available",
		"system_info": jsonWrap{
			"status":      "available",
			"environment": a.config.env,
			"version":     a.config.version,
		},
	}

	err := a.writeJSON(w, http.StatusOK, res, nil)
	if err != nil {
		a.logger.Error(err.Error())
	}
}
