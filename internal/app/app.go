package app

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"

	"temp/internal/config" // TODO: Change 'temp' to required name
	"temp/internal/server" // TODO: Change 'temp' to required name
)

type App struct {
	log    *slog.Logger
	server *server.Server
	cfg    *config.Config
	// Компоненты приложения
}

func New(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *slog.Logger, isShuttingDown *atomic.Bool) (*App, error) {
	// Инициализация компонентов

	// Инициализация сервера
	srv := server.NewServer(ctx, log, isShuttingDown)

	return &App{
		log:    log,
		server: srv,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	a.log.Info("Запуск HTTP сервера по адресу '" + a.cfg.URL + ":" + a.cfg.Port + "'...")
	return a.server.Start(a.cfg.Env, a.cfg.URL+":"+a.cfg.Port)
}

func (a *App) ShutDown(shutDownCtx context.Context) error {
	if a == nil {
		return errors.New("App instance is nil")
	}

	err := errors.Join(
		a.server.ShutDown(shutDownCtx),
	)
	return err
}
