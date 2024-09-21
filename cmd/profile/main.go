package main

import (
	"log/slog"
	"os"
	"os/signal"
	"profile/internal/app"
	"profile/internal/config"
	"profile/internal/database"
	"profile/internal/pkg/logger/handlers/slogpretty"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	postgres := database.NewPostgres(log, cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password,
		cfg.DB.Name)
	db := postgres.MustConnect()
	postgres.Ping()

	application := app.New(log, cfg.GRPC.Port, db)

	go func() { application.GRPC.MustRun() }()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))
	application.GRPC.Stop()
	postgres.Close()
	log.Info("gracefully stopped")
}

func setupLogger(env config.EnvType) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = setupPrettySlog()
	case config.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
