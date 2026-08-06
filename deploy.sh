#!/bin/bash
# ChengetAi Platform - Quick Deployment Script
# Usage: ./deploy.sh [environment]
# Example: ./deploy.sh production

set -e

ENVIRONMENT=${1:-development}

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          ChengetAi Platform - Deployment Script               ║"
echo "║                   Environment: $ENVIRONMENT                        ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Check prerequisites
echo "[1/8] Checking prerequisites..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi
if ! command -v docker compose &> /dev/null; then
    echo "❌ Docker Compose is not installed"
    exit 1
fi
echo "✓ Docker and Docker Compose are installed"
echo ""

# Check environment file
echo "[2/8] Checking environment configuration..."
if [ ! -f .env ]; then
    echo "⚠ .env file not found. Creating from .env.example..."
    cp .env.example .env
    echo "✓ .env created. Please edit it with your configuration:"
    echo "   nano .env"
    exit 1
fi
echo "✓ .env file exists"
echo ""

# Verify Docker daemon
echo "[3/8] Starting Docker daemon..."
if ! docker ps &> /dev/null; then
    echo "⚠ Docker daemon is not running. Attempting to start..."
    sudo systemctl start docker 2>/dev/null || {
        echo "❌ Could not start Docker daemon. Please start it manually:"
        echo "   sudo systemctl start docker"
        exit 1
    }
fi
echo "✓ Docker daemon is running"
echo ""

# Build images
echo "[4/8] Building service images..."
docker compose build --quiet
echo "✓ All services built successfully"
echo ""

# Start infrastructure
echo "[5/8] Starting infrastructure services..."
docker compose up -d postgres redis nats minio typesense
echo "✓ Infrastructure services started"
echo ""

# Wait for database
echo "[6/8] Waiting for database to be ready..."
max_attempts=30
attempt=1
while ! docker compose exec postgres psql -U chengetai -d chengetai -c "SELECT 1" &> /dev/null; do
    if [ $attempt -ge $max_attempts ]; then
        echo "❌ Database failed to become ready after ${max_attempts} attempts"
        exit 1
    fi
    echo "  Waiting... ($attempt/$max_attempts)"
    sleep 2
    ((attempt++))
done
echo "✓ Database is ready"
echo ""

# Start all services
echo "[7/8] Starting all microservices..."
docker compose up -d
echo "✓ All services started"
echo ""

# Verify deployment
echo "[8/8] Verifying deployment..."
echo ""
echo "Checking service health..."

services=(
    "api-gateway:8000"
    "auth-service:8001"
    "user-service:8002"
    "school-service:8003"
    "wallet-service:8004"
    "payment-service:8005"
    "notification-service:8006"
    "search-service:8008"
    "ai-service:8009"
    "quiz-service:8010"
    "recommendations-service:8011"
    "analytics-service:8012"
)

failed=0
for service_info in "${services[@]}"; do
    IFS=':' read -r service port <<< "$service_info"
    if docker compose exec "$service" curl -sf http://localhost:"$port"/health &> /dev/null; then
        echo "✓ $service (port $port)"
    else
        echo "❌ $service (port $port) - NOT READY"
        ((failed++))
    fi
done

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                  Deployment Complete!                          ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

if [ $failed -eq 0 ]; then
    echo "✓ All services are healthy and ready"
    echo ""
    echo "Access your platform at:"
    echo "  API Gateway:      http://localhost:8000"
    echo "  Auth Service:     http://localhost:8001"
    echo "  Search Service:   http://localhost:8008"
    echo "  AI Service:       http://localhost:8009"
    echo "  Quiz Service:     http://localhost:8010"
    echo "  Recommendations:  http://localhost:8011"
    echo "  Analytics:        http://localhost:8012"
    echo ""
    echo "Monitoring:"
    echo "  Grafana:          http://localhost:3000 (admin/admin)"
    echo "  Prometheus:       http://localhost:9090"
    echo "  MinIO:            http://localhost:9001 (minioadmin/minioadmin)"
    echo ""
    echo "Useful commands:"
    echo "  docker compose logs -f [service]    # View service logs"
    echo "  docker compose ps                   # Check service status"
    echo "  docker compose down                 # Stop all services"
    echo "  make health-check                   # Run health checks"
    echo ""
else
    echo "⚠ $failed service(s) are not yet healthy"
    echo ""
    echo "Check logs with:"
    echo "  docker compose logs -f"
    echo ""
    echo "Try again in a few seconds once services fully initialize"
fi
