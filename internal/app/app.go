package app

import (
	grpcapp "SSO/internal/app/grpc"
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
	// TODO: implement storage

	// TODO: init auth service (auth)

	grpcApp := grpcapp.New(log, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
