package app

import (
	"database/sql"
	"log/slog"
	grpcapp "profile/internal/app/grpc"
	"profile/internal/services"
	"profile/internal/storages"
	"profile/internal/transport"
)

type App struct {
	GRPC     *grpcapp.App
	database *sql.DB
}

func New(log *slog.Logger, port int, db *sql.DB) *App {
	storage := storages.NewProfileStorage(log, db)
	service := services.NewProfileService(log, storage)
	server := transport.NewServer(log, service)

	grpcApp := grpcapp.New(log, port, server)

	return &App{
		GRPC: grpcApp,
	}
}
