package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() //nolint:errcheck

	logger.Info("worker stub started")
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logger.Info("worker stub stopped", zap.String("signal", sig.String()))
}
