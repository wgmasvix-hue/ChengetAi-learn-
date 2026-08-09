package unit

import (
	"testing"

	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/config"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("APP_NAME", "chengetai")
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("REDIS_URL", "redis://localhost:6379")
	t.Setenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")
	t.Setenv("RATE_LIMIT_RPS", "12")
	t.Setenv("RATE_LIMIT_BURST", "24")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}
	if cfg.AppName != "chengetai" || cfg.AppPort != "9090" {
		t.Fatalf("unexpected config values: %+v", cfg)
	}
	if cfg.RateLimitRPS != 12 || cfg.RateLimitBurst != 24 {
		t.Fatalf("unexpected rate limit values: %+v", cfg)
	}
}

func TestLoadConfigRequiresCoreValues(t *testing.T) {
	t.Setenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz123456")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected missing env validation error")
	}
}
