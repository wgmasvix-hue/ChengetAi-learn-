package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppName        string
	AppEnv         string
	AppPort        string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	AIProvider     string
	AIAPIKey       string
	AIModel        string
	DSpaceBaseURL  string
	DSpaceAPIKey   string
	RateLimitRPS   int
	RateLimitBurst int
}

func Load() (Config, error) {
	cfg := Config{
		AppName:        getEnv("APP_NAME", "chengetai-learn-api"),
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8080"),
		DatabaseURL:    strings.TrimSpace(os.Getenv("DATABASE_URL")),
		RedisURL:       strings.TrimSpace(os.Getenv("REDIS_URL")),
		JWTSecret:      strings.TrimSpace(os.Getenv("JWT_SECRET")),
		AIProvider:     strings.TrimSpace(os.Getenv("AI_PROVIDER")),
		AIAPIKey:       strings.TrimSpace(os.Getenv("AI_API_KEY")),
		AIModel:        strings.TrimSpace(os.Getenv("AI_MODEL")),
		DSpaceBaseURL:  strings.TrimSpace(os.Getenv("DSPACE_BASE_URL")),
		DSpaceAPIKey:   strings.TrimSpace(os.Getenv("DSPACE_API_KEY")),
		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 20),
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.RedisURL == "" {
		missing = append(missing, "REDIS_URL")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if cfg.RateLimitRPS <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_RPS must be greater than zero")
	}
	if cfg.RateLimitBurst <= 0 {
		return Config{}, fmt.Errorf("RATE_LIMIT_BURST must be greater than zero")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
