package http

import (
	"net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/api-gateway/internal/handlers"
	"go.uber.org/zap"
)

func SetupRoutes(srv *Server, logger *zap.SugaredLogger) http.Handler {
	mux := http.NewServeMux()

	// Health checks
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/ready", handlers.ReadyHandler)

	// API v1
	mux.HandleFunc("/api/v1/ping", handlers.PingHandler)

	return mux
}
