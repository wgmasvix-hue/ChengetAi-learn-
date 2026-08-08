#!/bin/sh

# ChengetAi Render Entrypoint Script
# Handles service startup with proper error handling and health checks

set -e

# Verify required environment variables
if [ -z "$SERVICE_NAME" ]; then
    echo "ERROR: SERVICE_NAME not set"
    exit 1
fi

if [ -z "$SERVICE_PORT" ]; then
    echo "ERROR: SERVICE_PORT not set"
    exit 1
fi

# Log startup info
echo "============================================"
echo "ChengetAi - $SERVICE_NAME"
echo "============================================"
echo "Service: $SERVICE_NAME"
echo "Port: $SERVICE_PORT"
echo "Environment: ${ENVIRONMENT:-development}"
echo "Log Level: ${LOG_LEVEL:-info}"
echo ""

# For services that need database, verify connection
if [ -n "$DB_HOST" ] && [ -n "$DB_USER" ]; then
    echo "Checking database connectivity..."
    DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:${DB_PORT:-5432}/$DB_NAME?sslmode=${DB_SSL_MODE:-disable}"

    # Try to connect (retry up to 10 times)
    max_attempts=10
    attempt=1
    while [ $attempt -le $max_attempts ]; do
        if psql "$DB_URL" -c "SELECT 1" > /dev/null 2>&1; then
            echo "✓ Database connection successful"
            break
        fi
        if [ $attempt -lt $max_attempts ]; then
            echo "  Waiting for database... ($attempt/$max_attempts)"
            sleep 2
        else
            echo "✗ Database connection failed after $max_attempts attempts"
            echo "  Continuing anyway (database may be initializing)"
        fi
        attempt=$((attempt + 1))
    done
    echo ""
fi

# For services with Redis, verify connection
if [ -n "$REDIS_URL" ]; then
    echo "Checking Redis connectivity..."
    if redis-cli -u "$REDIS_URL" ping > /dev/null 2>&1; then
        echo "✓ Redis connection successful"
    else
        echo "⚠ Redis connection failed (service may not have Redis support)"
    fi
    echo ""
fi

# Validate critical secrets in production
if [ "$ENVIRONMENT" = "production" ]; then
    echo "Validating production configuration..."

    if [ -z "$JWT_SECRET" ]; then
        echo "ERROR: JWT_SECRET not set in production"
        exit 1
    fi

    if [ -n "$DB_HOST" ] && [ -z "$DB_PASSWORD" ]; then
        echo "ERROR: DB_PASSWORD not set in production"
        exit 1
    fi

    if [ ${#JWT_SECRET} -lt 32 ]; then
        echo "ERROR: JWT_SECRET too short (minimum 32 characters)"
        exit 1
    fi

    echo "✓ Production configuration validated"
    echo ""
fi

# Start the service
echo "Starting $SERVICE_NAME service..."
echo "Listening on port $SERVICE_PORT"
echo ""

# Execute the compiled service binary
exec /app/service
