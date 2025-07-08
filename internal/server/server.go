package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"

	sl "temp/internal/libs/logger" // TODO: Change 'temp' to required name

	"github.com/gorilla/csrf"
)

type Server struct {
	ctx            context.Context
	router         *http.ServeMux
	log            *slog.Logger
	isShuttingDown *atomic.Bool
	server         http.Server
}

func NewServer(ctx context.Context, log *slog.Logger, isShuttingDown *atomic.Bool) *Server {
	s := &Server{
		ctx:            ctx,
		router:         http.NewServeMux(),
		log:            log,
		isShuttingDown: isShuttingDown,
	}
	// Подключение функций к ручкам

	s.router.HandleFunc("GET /healthz", s.handleHealthz)

	return s
}

func (s *Server) Start(env, addr string) error {
	if env != "prod" {
		csrf.Secure(false)
		s.log.Debug("Not a prod")
	}
	// csrfProt := csrf.Protect([]byte("32-byte-long-auth-key"))
	s.server = http.Server{
		Addr:    addr,
		Handler: s.router,
		BaseContext: func(_ net.Listener) context.Context {
			return s.ctx
		},
	}

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.log.Error("Server failed to start", sl.Err(err))
		return err
	}
	return nil
}

// Функции для обработки запросов

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if s.isShuttingDown.Load() {
		http.Error(w, "Shutting down", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintln(w, "OK")
}

func (s *Server) ShutDown(shutDownCtx context.Context) error {
	return s.server.Shutdown(shutDownCtx)
}
