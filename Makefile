.PHONY: help build test clean docker-build docker-up docker-down docker-logs dev install-deps lint format

help:
	@echo "ChengetAi Platform - Development Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make build         - Build all services"
	@echo "  make test          - Run all tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-up     - Start all services with Docker Compose"
	@echo "  make docker-down   - Stop all services"
	@echo "  make docker-logs   - View Docker Compose logs"
	@echo "  make dev           - Start development environment"
	@echo "  make install-deps  - Install Go dependencies"
	@echo "  make lint          - Run linters"
	@echo "  make format        - Format Go code"

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

build:
	@echo "Building ChengetAi services..."
	@mkdir -p bin
	@for service in api-gateway auth-service user-service school-service wallet-service payment-service notification-service analytics-service search-service; do \
		echo "Building $$service..."; \
		cd apps/$$service && make build && cd ../..; \
	done
	@echo "✓ All services built"

test:
	@echo "Running tests..."
	@go test -v -cover ./apps/...
	@echo "✓ All tests passed"

lint:
	@echo "Linting code..."
	@go vet ./...
	@echo "✓ Lint passed"

format:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

clean:
	@echo "Cleaning build artifacts..."
	@find apps -name bin -type d -exec rm -rf {} + 2>/dev/null || true
	@go clean
	@echo "✓ Clean complete"

docker-build:
	@echo "Building Docker images..."
	@docker-compose build
	@echo "✓ Docker images built"

docker-up:
	@echo "Starting ChengetAi platform..."
	@docker-compose up -d
	@echo "✓ Platform started"
	@echo ""
	@echo "Services available at:"
	@echo "  API Gateway:      http://localhost:8000"
	@echo "  Auth Service:     http://localhost:8001"
	@echo "  User Service:     http://localhost:8002"
	@echo "  School Service:   http://localhost:8003"
	@echo "  Wallet Service:   http://localhost:8004"
	@echo "  Payment Service:  http://localhost:8005"
	@echo "  Notification Service: http://localhost:8006"
	@echo "  Analytics Service: http://localhost:8007"
	@echo "  Search Service:   http://localhost:8008"
	@echo ""
	@echo "Monitoring:"
	@echo "  Prometheus:       http://localhost:9090"
	@echo "  Grafana:          http://localhost:3000 (admin/admin)"
	@echo "  MinIO:            http://localhost:9001 (minioadmin/minioadmin)"
	@echo ""
	@docker-compose ps

docker-down:
	@echo "Stopping ChengetAi platform..."
	@docker-compose down
	@echo "✓ Platform stopped"

docker-logs:
	@docker-compose logs -f

dev: docker-up
	@echo "✓ Development environment ready"
	@echo "Run 'make docker-logs' to see logs"

health-check:
	@echo "Checking service health..."
	@for port in 8000 8001 8002 8003 8004 8005 8006 8007 8008; do \
		echo -n "Port $$port: "; \
		curl -s -o /dev/null -w "%{http_code}" http://localhost:$$port/health || echo "OFFLINE"; \
		echo ""; \
	done

migrate:
	@echo "Running database migrations..."
	@docker-compose exec postgres psql -U chengetai -d chengetai -f /docker-entrypoint-initdb.d/migrations.sql
	@echo "✓ Migrations complete"

logs-%:
	@docker-compose logs -f $*

ps:
	@docker-compose ps

env-check:
	@echo "Checking environment..."
	@[ -f .env ] && echo "✓ .env file exists" || echo "✗ .env file missing - copy from .env.example"
	@command -v docker >/dev/null 2>&1 && echo "✓ Docker installed" || echo "✗ Docker not installed"
	@command -v docker-compose >/dev/null 2>&1 && echo "✓ Docker Compose installed" || echo "✗ Docker Compose not installed"
	@command -v go >/dev/null 2>&1 && echo "✓ Go installed" || echo "✗ Go not installed"

init: env-check install-deps docker-build
	@echo "✓ ChengetAi environment initialized"
	@echo "Run 'make dev' to start the platform"

.DEFAULT_GOAL := help
