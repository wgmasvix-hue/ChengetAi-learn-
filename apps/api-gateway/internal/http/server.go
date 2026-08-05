package http

import (
	"net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/api-gateway/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	config *config.Config
	logger *zap.SugaredLogger
}

func NewServer(cfg *config.Config, logger *zap.SugaredLogger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

// ListenAndServe starts the HTTP server
func (s *Server) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx interface{}) error {
	return nil
}
