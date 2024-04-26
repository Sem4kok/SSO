package main

import (
	"SSO/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// project layer's schema
// --------------------------------
// 1) Transport layer.
// Transport layer gets request and
// communicate with Service layer
//
// 2) Service layer. (auth, permission, userinfo)
// Service layer implement business-logic
// communicate with Data layer
//
// 3) Data layer.
// Data layer communicate with data (include storage)
// return response's to Service layer

func main() {
	cfg := config.Load()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))
	// TODO: initialize logger

	// TODO: Initialize application(app)

	// TODO: start gRPC-server app
}

// setupLogger will choose logger setup
// choose will be based on env
// implemented by go built-in library slog
// log.Debug() <- 1   level's of logger messages
// log.Info()  <- 2	  envLocal logger has (1) level
// log.Warn()  <- 3   writes for debug - human-oriented format.
// log.Error() <- 4	  envDev has (1) level writes JSON format.
// envProd has (2) level writes human-oriented format
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
