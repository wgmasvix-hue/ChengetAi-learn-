package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// PostgresConfig holds PostgreSQL configuration
type PostgresConfig struct {
	URL            string
	MaxConnections int
	MinConnections int
	ConnTimeout    time.Duration
}

// NewPostgresConnection creates a new PostgreSQL connection pool
func NewPostgresConnection(ctx context.Context, cfg PostgresConfig, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	if cfg.MaxConnections == 0 {
		cfg.MaxConnections = 25
	}
	if cfg.MinConnections == 0 {
		cfg.MinConnections = 5
	}
	if cfg.ConnTimeout == 0 {
		cfg.ConnTimeout = 10 * time.Second
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MinConnections)
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnTimeout

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("creating database pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	logger.Info("PostgreSQL connection pool created",
		"max_connections", cfg.MaxConnections,
		"min_connections", cfg.MinConnections,
	)

	return pool, nil
}

// HealthCheck checks if the database is healthy
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Release()

	return conn.Ping(ctx)
}
