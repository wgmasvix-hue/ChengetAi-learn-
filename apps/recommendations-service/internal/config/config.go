package config

import (
	"os"
	"strconv"

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

	// Recommendations
	DefaultTopK         int
	DefaultMinSimilarity float32
	DefaultStrategy     string // collaborative, content_based, hybrid

	// Logging
	LogLevel string
}

func Load() *Config {
	return &Config{
		ServiceName:         getEnv("SERVICE_NAME", "recommendations-service"),
		Port:                getEnvInt("SERVICE_PORT", 8011),
		Environment:         getEnv("ENVIRONMENT", "development"),
		DatabaseURL:         buildDatabaseURL(),
		DatabaseMaxConns:    getEnvInt("DB_MAX_CONNS", 25),
		DatabaseMinConns:    getEnvInt("DB_MIN_CONNS", 5),
		DefaultTopK:         getEnvInt("DEFAULT_TOP_K", 10),
		DefaultMinSimilarity: getEnvFloat("DEFAULT_MIN_SIMILARITY", 0.5),
		DefaultStrategy:     getEnv("DEFAULT_STRATEGY", "hybrid"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
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

func getEnvFloat(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatVal)
		}
	}
	return defaultValue
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
