# ChengetAi Learn

ChengetAi Learn is Phase One of a Zimbabwe-focused learning platform that combines a Go API, PostgreSQL, Redis, and a lightweight web frontend for learner, teacher, and parent onboarding.

## Prerequisites

- Go 1.22
- Docker and Docker Compose

## Quick start

```bash
cp .env.example .env
docker compose up -d
```

The API is available at `http://localhost:8080` and the frontend at `http://localhost:8081`.

`DATABASE_URL` is used for host-based development and `DOCKER_DATABASE_URL` is used by the API container inside Docker Compose.

## API endpoints

- `GET /health`
- `GET /ready`
- `GET /api/v1/version`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`
- `GET /api/v1/users/me`
- `GET /api/v1/learners/me`
- `GET /api/v1/teachers/me`

## Development commands

```bash
make dev
make build
make test
make lint
make migrate
make migrate-down
make seed
make docker-up
make docker-down
```

## Architecture overview

- `cmd/` contains application entrypoints for the API, worker stub, and migrations.
- `internal/` contains configuration, infrastructure, middleware, and domain packages.
- `migrations/` contains PostgreSQL schema migrations.
- `web/` contains a framework-free mobile-first frontend.
- `docs/openapi.yaml` documents the initial HTTP API.

## Authentication

JWTs are signed with HS256 and stored in Redis as active session references so logout can invalidate tokens.
