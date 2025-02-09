package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// server creates a new server from the application's configuration parameters
// and middleware.
func (app *app) server() error {
	// handler is the serve mux, wrapped with appropriate middleware.
	//var handler http.Handler = app.recoverPanic(
	//	app.enableCORS(
	//		app.rateLimit(
	//			app.
	//				routes(),
	//		),
	//	),
	//)
	var handler http.Handler = app.enableCORS(
		app.routes(),
	)

	// handler = app.enableCORS(app.rateLimit(app.routes()))
	handler = app.enableCORS(app.routes())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Info("starting server", slog.Int("port", app.config.port))

	return srv.ListenAndServe()
}
