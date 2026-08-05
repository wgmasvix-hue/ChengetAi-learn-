package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/recommendations-service/internal/config"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/recommendations-service/internal/handlers"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/recommendations-service/internal/services"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/database"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/recommendations"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logging
	logger := config.SetupLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("Starting Recommendations Service",
		"port", cfg.Port,
		"environment", cfg.Environment,
	)

	// Setup database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.NewPostgresConnection(ctx, database.PostgresConfig{
		URL:            cfg.DatabaseURL,
		MaxConnections: cfg.DatabaseMaxConns,
		MinConnections: cfg.DatabaseMinConns,
		ConnTimeout:    10 * time.Second,
	}, logger)
	if err != nil {
		logger.Fatalw("Failed to connect to database", "error", err)
	}
	defer dbPool.Close()

	// Initialize recommendation engine
	recEngine := recommendations.NewEngine(dbPool, logger)

	// Initialize service
	recService := services.NewRecommendationService(recEngine, dbPool, logger)

	// Setup routes
	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := database.HealthCheck(context.Background(), dbPool); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ready"}`)
	})

	// Recommendation endpoints
	recHandler := handlers.NewRecommendationHandler(recService, logger)
	mux.HandleFunc("GET /recommendations/{userId}", recHandler.GetRecommendations)
	mux.HandleFunc("POST /recommendations/generate", recHandler.GenerateRecommendations)
	mux.HandleFunc("GET /metrics/{userId}", recHandler.GetMetrics)
	mux.HandleFunc("POST /feedback", recHandler.SubmitFeedback)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Infow("HTTP server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorw("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	sig := <-sigChan
	logger.Infow("Shutdown signal received", "signal", sig.String())

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorw("Shutdown error", "error", err)
	}

	logger.Info("Recommendations Service stopped")
}
