// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

// // routeHealthcheck is a default route for all HTTP services. Hit this
// // endpoint to check if the service is down or to check on degradation status
// func routeHealthcheck(mux *http.ServeMux) {
// 	mux.HandleFunc(
// 		"/healthcheck", func(w http.ResponseWriter, r *http.Request) {
// 			status := jsonWrap{
// 				"status": "healthy",
// 				"uptime": "",
// 			}
// 		},
// 	)
// }
