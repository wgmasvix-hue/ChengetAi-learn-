# API Gateway Service

The API Gateway is the entry point for all client requests to the ChengetAi Platform. It handles:
- Request routing to appropriate microservices
- Authentication and authorization validation
- Rate limiting and DDoS protection
- Request/response logging and monitoring
- API versioning

## Architecture

```
Client Requests
    ↓
Caddy (Reverse Proxy)
    ↓
API Gateway (Port 8000)
    ├── Route to Auth Service (Port 8001)
    ├── Route to User Service (Port 8002)
    ├── Route to School Service (Port 8003)
    ├── Route to Wallet Service (Port 8004)
    ├── Route to Payment Service (Port 8005)
    ├── Route to Notification Service (Port 8006)
    ├── Route to Analytics Service (Port 8007)
    └── Route to Search Service (Port 8008)
```

## Building

```bash
make build
```

## Running Locally

```bash
# Set environment variables
export SERVICE_NAME=api-gateway
export SERVICE_PORT=8000
export DB_HOST=localhost
export REDIS_URL=redis://localhost:6379
export NATS_URL=nats://localhost:4222

# Run the service
make run
```

## Docker

Build the Docker image:
```bash
make docker-build
```

Run in Docker:
```bash
make docker-run
```

## Testing

```bash
make test
```

## API Endpoints

### Health Checks
- `GET /health` - Service health status
- `GET /ready` - Service readiness status

### API v1
- `GET /api/v1/ping` - Test endpoint

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVICE_NAME` | Service identifier | api-gateway |
| `SERVICE_PORT` | HTTP port | 8000 |
| `ENVIRONMENT` | Environment (dev/prod) | development |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_NAME` | Database name | chengetai |
| `DB_USER` | Database user | chengetai |
| `DB_PASSWORD` | Database password | (required in production) |
| `REDIS_URL` | Redis connection URL | redis://localhost:6379 |
| `NATS_URL` | NATS server URL | nats://localhost:4222 |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | info |

## Middleware

The API Gateway uses the following middleware stack:

1. **RequestID**: Adds unique request ID for tracing
2. **RequestLogger**: Logs all incoming requests
3. **Recovery**: Recovers from panics
4. **CORS**: Handles Cross-Origin Resource Sharing

## Next Steps

- [ ] Implement service discovery
- [ ] Add rate limiting
- [ ] Implement circuit breakers
- [ ] Add request validation
- [ ] Implement authentication middleware
- [ ] Add API documentation (OpenAPI/Swagger)
