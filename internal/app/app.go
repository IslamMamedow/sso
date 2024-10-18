package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	GRPCport int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO инициализировать хранилище

	// TODO инициализировать auth сервер

	grpcServer := grpcapp.New(log, GRPCport)

	return &App{
		GRPCServer: grpcServer,
	}
}
