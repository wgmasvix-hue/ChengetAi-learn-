package http

import (
	stdhttp "net/http"
	"time"

	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/config"
)

func NewServer(cfg config.Config, handler stdhttp.Handler) *stdhttp.Server {
	return &stdhttp.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
