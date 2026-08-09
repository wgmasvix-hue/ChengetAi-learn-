# ChengetAi Learn

> Learn smarter with ChengetAi

ChengetAi Learn is a Zimbabwe-focused, mobile-first education platform. It is being developed at **[learn.dare.co.zw](https://learn.dare.co.zw)** as part of the DARE (Digital Access and Resources for Education) initiative.

Phase One builds the production-ready foundation: Go API → PostgreSQL → Redis → JWT auth → Docker → CI.

## Prerequisites

- Go 1.22+
- Docker and Docker Compose

## Quick start

```bash
cp .env.example .env
# Edit .env and set a strong JWT_SECRET (minimum 32 characters)
docker compose up -d
```

| Service  | Local URL                    |
|----------|------------------------------|
| Frontend | http://localhost:8081        |
| API      | http://localhost:8080        |
| Health   | http://localhost:8080/health |

`DATABASE_URL` is used for host-based development; `DOCKER_DATABASE_URL` is used by the API container inside Docker Compose.

## Production (learn.dare.co.zw)

In production the site is served over HTTPS by Caddy, which handles TLS automatically via Let's Encrypt.

Set these in your production `.env`:

```
APP_ENV=production
APP_URL=https://learn.dare.co.zw
CORS_ALLOWED_ORIGINS=https://learn.dare.co.zw
CSP_CONNECT_SRC='self' https://learn.dare.co.zw
```

Then uncomment the `caddy` service block in `docker-compose.yml` before deploying.

## API endpoints

| Method | Path                    | Auth     | Description              |
|--------|-------------------------|----------|--------------------------|
| GET    | /health                 | —        | Liveness probe           |
| GET    | /ready                  | —        | Readiness (postgres+redis)|
| GET    | /api/v1/version         | —        | API version              |
| POST   | /api/v1/auth/register   | —        | Register learner/teacher/parent |
| POST   | /api/v1/auth/login      | —        | Login                    |
| POST   | /api/v1/auth/logout     | ****** Logout (invalidates token)|
| GET    | /api/v1/auth/me         | ****** Current user             |
| GET    | /api/v1/users/me        | ****** User profile             |
| GET    | /api/v1/learners/me     | ****** Learner profile          |
| GET    | /api/v1/teachers/me     | ****** Teacher profile          |

Full OpenAPI specification: [`docs/openapi.yaml`](docs/openapi.yaml)

## Development commands

```bash
make dev          # run API locally (go run)
make build        # compile binary to bin/api
make test         # run all tests
make lint         # go vet + gofmt check
make migrate      # apply database migrations
make migrate-down # roll back last migration
make seed         # seed development data
make docker-up    # docker compose up -d
make docker-down  # docker compose down
```

## Architecture

```
Frontend (web/)
     ↓
Go API (cmd/api/)
     ↓
Middleware (logging · auth · CORS · rate limit)
     ↓
Handlers → Services → Repositories
     ↓            ↓
  PostgreSQL     Redis
```

Package layout:

- `cmd/` — entrypoints: `api`, `worker`, `migrate`
- `internal/config/` — environment-variable config
- `internal/database/` — pgxpool connection
- `internal/redis/` — Redis client
- `internal/http/` — router and server
- `internal/middleware/` — logging, auth, security
- `internal/auth/` — register, login, logout, me
- `internal/users/` — user profile
- `internal/learners/` — learner profile
- `internal/teachers/` — teacher profile
- `internal/health/` — /health and /ready
- `internal/ai/` — AI provider interface (Phase 4)
- `internal/knowledge/` — knowledge repository interface (Phase 5)
- `migrations/` — versioned SQL migrations
- `web/` — mobile-first vanilla JS/HTML/CSS frontend
- `docs/openapi.yaml` — OpenAPI 3.0 specification

## Authentication

Passwords are hashed with **Argon2id**. JWTs are signed with HS256 and stored as session references in Redis so logout immediately invalidates the token.

Public registration supports roles: `learner`, `teacher`, `parent`. Administrative roles (`school_admin`, `platform_admin`) cannot be self-assigned.
