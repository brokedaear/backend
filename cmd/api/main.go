package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"backend.brokedaear.com/utils/prettylog"
)

func main() {
	var err error

	envFlag := flag.Int("environment", 0, "application environment 1-3")
	flag.Parse()

	if *envFlag < 0 || *envFlag > 2 {
		panic("invalid application environment")
	}

	config := &config{
		port:    7402,
		env:     Environment(*envFlag),
		version: "0.1",
	}

	slogHandlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if config.env == ENV_DEVELOPMENT || config.env == ENV_STAGING {
		slogHandlerOptions.AddSource = true
		slogHandlerOptions.Level = slog.LevelDebug
	}

	logger := slog.New(prettylog.NewHandler(slogHandlerOptions))

	logger.Info(fmt.Sprintf("application environment set to %d", config.env))

	app := &app{
		config: config,
		logger: logger,
	}

	err = app.server()

	logger.Error(err.Error())
	os.Exit(1)
}
