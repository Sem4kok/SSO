package grpcapp

import (
	"SSO/internal/grpc/auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New is constructor for app
func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	auth.Register(gRPCServer)

	return &App{
		log:        log,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

// Run starts gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	// give attributes to app logger
	// that for each log msg get info
	// where error was appeared
	log := a.log.With(slog.String("op", op))

	// creating port using fmt.Sprintf
	port := fmt.Sprintf(":%d", a.port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is listening", slog.String("address", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info(
		"stopping gRPC server", slog.Int("port", a.port))

	// needed to gracefully stop the server
	a.gRPCServer.GracefulStop()
}
