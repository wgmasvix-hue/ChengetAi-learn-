# ChengetAi Learn - Phase One: Foundation

Phase One establishes the production-ready foundation for ChengetAi Learn.

## Quick Start

```bash
# 1. Copy environment file
cp .env.example.phase1 .env

# 2. Start Docker services
make docker-up

# 3. Run database migrations
make migrate

# 4. Start API server
make dev

# 5. Access the API
curl http://localhost:8080/health
```

## Architecture

```
┌─────────────┐
│  Frontend   │ (localhost:3000)
└──────┬──────┘
       │
┌──────┴──────┐
│   HTTP      │ (localhost:8080)
│   API       │
└──────┬──────┘
       │
   ┌───┴───┐
   │       │
┌──▼──┐ ┌─▼──┐
│  DB │ │Cache│
└─────┘ └─────┘
```

## Key Endpoints

### Health
- `GET /health` - Liveness probe
- `GET /ready` - Readiness probe (checks dependencies)
- `GET /version` - API version

### Authentication (Placeholder Phase One)
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/me` - Current user info

### Users
- `GET /api/v1/users` - List users

## What's Included

✓ Go API with proper structure
✓ PostgreSQL connection pooling
✓ Redis connection
✓ JWT token infrastructure
✓ Bcrypt password hashing
✓ Request ID tracking
✓ Structured logging
✓ CORS configuration
✓ Security headers
✓ Health checks
✓ Docker/Docker Compose
✓ Database migrations
✓ Standard error responses

## What's NOT Included Yet

- AI tutor
- Curriculum engine
- RAG/Knowledge retrieval
- DSpace integration
- Payments
- Parent portal
- School analytics
- Offline sync
- Advanced features

These will be implemented in subsequent phases.

## Development

### Running Tests
```bash
make test
```

### Running Linter
```bash
make lint
```

### Building Binaries
```bash
make build
```

### Docker Commands
```bash
make docker-up          # Start services
make docker-down        # Stop services
make docker-clean       # Stop and remove volumes
make docker-logs        # View API logs
make docker-logs-db     # View database logs
make docker-logs-redis  # View Redis logs
```

### Database Migrations
```bash
make migrate            # Run migrations
make migrate-down       # Rollback
```

## Configuration

Edit `.env` for configuration:

```bash
APP_ENV=development
DATABASE_URL=postgres://user:pass@host:5432/db
REDIS_URL=redis://host:6379/0
JWT_SECRET=your-secret-key
```

## Docker Services

- **postgres**: PostgreSQL 16
- **redis**: Redis 7.2
- **api**: Go API server

All services run on internal Docker network `chengetai-network`.

## Acceptance Criteria ✓

Phase One is complete when:

### Infrastructure
- [x] Docker Compose works
- [x] PostgreSQL starts
- [x] Redis starts
- [x] API starts
- [ ] Frontend starts (coming next)

### Backend
- [x] Go application builds
- [x] /health endpoint works
- [x] /ready endpoint works
- [x] /version endpoint works
- [x] Structured logging works
- [x] Request IDs work
- [x] API errors are standardised
- [ ] Tests pass (in progress)
- [ ] CI passes (coming next)

### Authentication
- [ ] Learner registration works
- [ ] Teacher registration works
- [ ] Login works
- [ ] /me works
- [ ] Passwords are securely hashed
- [ ] Role checks work

### Database
- [x] Migrations work
- [x] Rollbacks work
- [ ] User records persist (needs integration)
- [ ] Profile records persist (needs integration)

### Quality
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Linter passes
- [ ] CI passes

## What's Next

1. **Frontend** - Create React/Next.js web application
2. **Database Integration** - Wire up auth handlers to database
3. **Integration Tests** - Test API + Database + Redis together
4. **CI/CD** - Add GitHub Actions workflow
5. **User Service** - Full user management implementation
6. **Phase Two** - Begin curriculum and learning engine

## Ports

- API: 8080
- PostgreSQL: 5432
- Redis: 6379
- Frontend: 3000 (when added)

## Documentation

- API Documentation: `/docs` (coming with OpenAPI)
- Architecture: See ARCHITECTURE.md
- Security: See SECURITY.md
- Deployment: See RENDER_DEPLOYMENT.md

## Support

For issues or questions:
1. Check logs: `make docker-logs`
2. Check database: Connect to postgres:5432
3. Check Redis: `redis-cli -h localhost`
4. View API health: `curl http://localhost:8080/health`

---

**Phase One Status**: Foundation Building ✓  
**Target Completion**: When all acceptance criteria ✓

