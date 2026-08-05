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

	// Claude API
	ClaudeAPIKey   string
	ClaudeModel    string
	ClaudeMaxTokens int

	// RAG
	RAGMaxContextLength int
	RAGMinSimilarity    float32
	RAGTopK             int

	// Quiz
	DefaultQuestionCount int
	DefaultQuestionTypes string // comma-separated

	// Logging
	LogLevel string
}

func Load() *Config {
	return &Config{
		ServiceName:         getEnv("SERVICE_NAME", "ai-service"),
		Port:                getEnvInt("SERVICE_PORT", 8009),
		Environment:         getEnv("ENVIRONMENT", "development"),
		DatabaseURL:         buildDatabaseURL(),
		DatabaseMaxConns:    getEnvInt("DB_MAX_CONNS", 25),
		DatabaseMinConns:    getEnvInt("DB_MIN_CONNS", 5),
		ClaudeAPIKey:        getEnv("CLAUDE_API_KEY", ""),
		ClaudeModel:         getEnv("CLAUDE_MODEL", "claude-opus-5"),
		ClaudeMaxTokens:     getEnvInt("CLAUDE_MAX_TOKENS", 4096),
		RAGMaxContextLength: getEnvInt("RAG_MAX_CONTEXT_LENGTH", 8192),
		RAGMinSimilarity:    getEnvFloat("RAG_MIN_SIMILARITY", 0.5),
		RAGTopK:             getEnvInt("RAG_TOP_K", 5),
		DefaultQuestionCount: getEnvInt("DEFAULT_QUESTION_COUNT", 10),
		DefaultQuestionTypes: getEnv("DEFAULT_QUESTION_TYPES", "multiple_choice,true_false,short_answer"),
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
