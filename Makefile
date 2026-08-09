.PHONY: dev build test lint migrate migrate-down seed docker-up docker-down docker-build

dev:
	go run ./cmd/api/...

build:
	mkdir -p bin && go build -o bin/api ./cmd/api/...

test:
	go test ./...

lint:
	go vet ./...
	gofmt -l .

migrate:
	go run ./cmd/migrate/... up

migrate-down:
	go run ./cmd/migrate/... down

seed:
	psql $$DATABASE_URL -f scripts/seed.sql

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build
