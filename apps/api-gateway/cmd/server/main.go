package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/api-gateway/internal/config"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/api-gateway/internal/http"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/api-gateway/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logging
	logger := config.SetupLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("Starting API Gateway",
		"port", cfg.Port,
		"environment", cfg.Environment,
	)

	// Setup HTTP server
	srv := http.NewServer(cfg, logger)
	router := http.SetupRoutes(srv, logger)

	// Add middleware
	handler := middleware.Chain(
		router,
		middleware.RequestID(),
		middleware.RequestLogger(logger),
		middleware.Recovery(),
		middleware.CORS(cfg),
	)

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		logger.Info("HTTP server listening", "addr", addr)
		if err := srv.ListenAndServe(addr, handler); err != nil {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	logger.Info("Shutdown signal received", "signal", sig.String())

	// Shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error", "error", err)
	}

	logger.Info("API Gateway stopped")
}
