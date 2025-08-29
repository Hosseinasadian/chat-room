package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"time"
)

type Config struct {
	Host               string        `koanf:"host"`
	Port               int           `koanf:"port"`
	Pattern            string        `koanf:"pattern"`
	ShutDownCtxTimeout time.Duration `koanf:"shut_down_ctx_timeout"`
}

type Handler interface {
	Routes() chi.Router
}

type Server struct {
	config     Config
	handler    Handler
	httpServer *http.Server
}

func New(config Config, handler Handler) Server {
	return Server{
		config:  config,
		handler: handler,
	}
}

func (s *Server) Serve() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // React dev server
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	r.Mount(s.config.Pattern, s.handler.Routes())

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler: r,
	}

	// blocking until shutdown
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(httpShutdownCtx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(httpShutdownCtx)
}
