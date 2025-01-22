package main

import "net/http"

func (app *app) routes() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	return router
}
