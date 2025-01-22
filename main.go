package main

import (
	"log"
	"os"
)

const LOGGER_PREFIX = ""

func main() {
	var err error

	logger := log.New(os.Stdout, LOGGER_PREFIX, log.Ldate|log.Ltime)

	config := config{
		port: 3453,
	}

	app := &app{
		config: config,
		logger: logger,
	}

	err = app.server()

	logger.Fatal(err)
}

// config is application configuration.
type config struct {

	// port access for host.
	port int

	// Runtime environment, either "development", "staging", or "production".
	env string

	// db is the database configuration config.
	db dbConfig

	// limiter is for limiter information for rate limiting.
	limiter struct {
		// rps is requests per second.
		rps float64

		// burst is how many bursts are allowed.
		burst int

		// enabled either disables or enables rate limited altogether.
		enabled bool
	}
}

type dbConfig struct {
	// Driver is the SQL driver, like postgreSQL..
	Driver string

	// DSN is Data Source Name.
	Dsn string

	// Database parameters. On the app layer, these are found in env variables.

	Name     string
	Username string
	Password string
	Host     string
	Port     string
	SslMode  string

	// Connection rate limiting logic.

	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
