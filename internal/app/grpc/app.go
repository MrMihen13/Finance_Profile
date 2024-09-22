package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"profile/internal/transport"
)

const opAppName string = "grpcapp.App"

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	server     *transport.Server
	port       int
}

func New(log *slog.Logger, port int, srv *transport.Server) *App {
	app := &App{
		log:        log,
		gRPCServer: grpc.NewServer(),
		server:     srv,
		port:       port,
	}

	app.server.Register(app.gRPCServer)

	return app
}

func (a App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a App) Run() error {
	const op = opAppName + ".Run"
	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a App) Stop() {
	const op = opAppName + ".Stop"
	a.log.With(slog.String("op", op)).Info("stopping grpc server")
	a.gRPCServer.GracefulStop()
}
