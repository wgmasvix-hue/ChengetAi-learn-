package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Application
	AppName string
	AppEnv  string
	AppPort string
	LogLevel string

	// Database
	DatabaseURL string
	DBMaxConns  int
	DBMinConns  int

	// Redis
	RedisURL string

	// Authentication
	JWTSecret string
	JWTExpiry string

	// AI Configuration (future use)
	AIProvider string
	AIAPIKey   string
	AIModel    string

	// Knowledge Repository (future use)
	DSpaceBaseURL string
	DSpaceAPIKey  string

	// Security
	CORSAllowedOrigins string
	RateLimitEnabled   bool

	// Environment Detection
	IsDevelopment bool
	IsProduction  bool
	IsTest        bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		AppName:            getEnv("APP_NAME", "chengetai-learn"),
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		RedisURL:           getEnv("REDIS_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiry:          getEnv("JWT_EXPIRY", "15m"),
		AIProvider:         getEnv("AI_PROVIDER", ""),
		AIAPIKey:           getEnv("AI_API_KEY", ""),
		AIModel:            getEnv("AI_MODEL", ""),
		DSpaceBaseURL:      getEnv("DSPACE_BASE_URL", ""),
		DSpaceAPIKey:       getEnv("DSPACE_API_KEY", ""),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080"),
		RateLimitEnabled:   getBoolEnv("RATE_LIMIT_ENABLED", true),
	}

	cfg.DBMaxConns = getIntEnv("DB_MAX_CONNS", 10)
	cfg.DBMinConns = getIntEnv("DB_MIN_CONNS", 2)

	// Set environment flags
	cfg.IsDevelopment = cfg.AppEnv == "development"
	cfg.IsProduction = cfg.AppEnv == "production"
	cfg.IsTest = cfg.AppEnv == "test"

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks required configuration
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}

	if c.IsProduction {
		if c.JWTSecret == "" || len(c.JWTSecret) < 32 {
			return fmt.Errorf("JWT_SECRET is required and must be at least 32 characters in production")
		}
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}
