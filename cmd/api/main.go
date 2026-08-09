package main

import (
	"context"
	"errors"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appauth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/auth"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/config"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/database"
	apphealth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/health"
	apphttp "github.com/wgmasvix-hue/ChengetAi-learn-/internal/http"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/learners"
	redispkg "github.com/wgmasvix-hue/ChengetAi-learn-/internal/redis"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/teachers"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/users"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() //nolint:errcheck

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer dbPool.Close()

	redisClient, err := redispkg.New(cfg.RedisURL)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer redisClient.Close() //nolint:errcheck

	authRepo := appauth.NewRepository(dbPool, redisClient)
	authService := appauth.NewService(authRepo, dbPool, cfg.JWTSecret)
	authHandler := appauth.NewHandler(authService)

	usersHandler := users.NewHandler(users.NewService(users.NewRepository(dbPool)))
	learnersHandler := learners.NewHandler(learners.NewService(learners.NewRepository(dbPool)))
	teachersHandler := teachers.NewHandler(teachers.NewService(teachers.NewRepository(dbPool)))
	healthHandler := apphealth.NewHandler(
		func(ctx context.Context) error { return database.Ping(ctx, dbPool) },
		func(ctx context.Context) error { return redispkg.Ping(ctx, redisClient) },
	)

	router := apphttp.NewRouter(cfg, logger, redisClient, authHandler, usersHandler, learnersHandler, teachersHandler, healthHandler)
	server := apphttp.NewServer(cfg, router)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting api server", zap.String("addr", server.Addr), zap.String("environment", cfg.AppEnv))
		if serveErr := server.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, stdhttp.ErrServerClosed) {
			errCh <- serveErr
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	case serveErr := <-errCh:
		logger.Fatal("api server failed", zap.Error(serveErr))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	}
}
