package app

import (
	grpcapp "SSO/internal/app/grpc"
	"SSO/internal/services/auth"
	"SSO/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

// New is constructor for app
func New(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// storage initialization
	storage, err := sqlite.Connect(storagePath)
	if err != nil {
		panic(err)
	}

	// initialize app (auth)
	// storage realize interfaces
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
