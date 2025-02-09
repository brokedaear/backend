package main

import "net/http"

// healthCheckHandler handles healthcheck requests. It uses writeJSON()
// to write system health data to the http.ResponseWriter stream.
func (app *app) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	app.logger.Info("received %s from %s", r.Method, r.RemoteAddr)

	res := jsonWrap{
		"status": "available",
		"system_info": map[string]any{
			"status":      "available",
			"environment": app.config.env,
			"version":     app.config.version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, res, nil)

	if err != nil {
		app.logger.Error(err.Error())
	}
}
