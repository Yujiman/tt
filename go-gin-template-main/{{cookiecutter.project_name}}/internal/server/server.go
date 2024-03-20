package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"{{cookiecutter.package_name}}/internal/config"
	handler "{{cookiecutter.package_name}}/internal/delivery/http"
)

type Server struct {
	server *http.Server
}

func NewServer(cfg *config.Config, handler *handler.Handler) (*Server, error) {
	httpHandler, err := handler.Init(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		server: &http.Server{
			Addr:           fmt.Sprintf(":%s", cfg.Service.Port),
			Handler:        httpHandler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 2 << 20, // 2MB
		},
	}, nil
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
