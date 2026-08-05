package config

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// Service
	ServiceName string
	Port        int
	Environment string

	// Database
	DatabaseURL      string
	DatabaseMaxConns int
	DatabaseMinConns int

	// Redis (for future use)
	RedisURL string

	// NATS (for future use)
	NatsURL string

	// JWT
	JWTSecret            string
	JWTAccessDuration    time.Duration
	JWTRefreshDuration   time.Duration

	// Logging
	LogLevel string
}

func Load() *Config {
	return &Config{
		ServiceName:        getEnv("SERVICE_NAME", "auth-service"),
		Port:               getEnvInt("SERVICE_PORT", 8001),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        buildDatabaseURL(),
		DatabaseMaxConns:   getEnvInt("DB_MAX_CONNS", 25),
		DatabaseMinConns:   getEnvInt("DB_MIN_CONNS", 5),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		NatsURL:            getEnv("NATS_URL", "nats://localhost:4222"),
		JWTSecret:          getEnv("JWT_SECRET", "dev_jwt_secret_change_in_prod"),
		JWTAccessDuration:  parseDuration(getEnv("JWT_ACCESS_DURATION", "24h")),
		JWTRefreshDuration: parseDuration(getEnv("JWT_REFRESH_DURATION", "168h")),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
	}
}

func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	name := getEnv("DB_NAME", "chengetai")
	user := getEnv("DB_USER", "chengetai")
	password := getEnv("DB_PASSWORD", "")

	if password != "" {
		return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable"
	}
	return "postgres://" + user + "@" + host + ":" + port + "/" + name + "?sslmode=disable"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return d
}

// SetupLogger configures the logger
func SetupLogger(logLevel string) *zap.SugaredLogger {
	level := zapcore.InfoLevel
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger.Sugar()
}
