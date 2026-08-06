# ChengetAi Platform - Quick Start Guide

## 🚀 One-Command Launch

On your server, clone and launch:

```bash
git clone https://github.com/wgmasvix-hue/ChengetAi-learn-.git
cd ChengetAi-learn-
git checkout claude/finish-lixm6n
./deploy.sh production
```

The `deploy.sh` script handles everything:
1. ✓ Checks Docker/Docker Compose installation
2. ✓ Creates .env file (edit it with your config)
3. ✓ Builds all 13 services
4. ✓ Starts infrastructure (PostgreSQL, Redis, etc.)
5. ✓ Waits for database to be ready
6. ✓ Launches all microservices
7. ✓ Verifies health checks

## 📊 What's Included

### Core Services (Phases 0-2)
- **API Gateway** (8000) - Request routing & load balancing
- **Auth Service** (8001) - JWT authentication
- **User Service** (8002) - User management
- **School Service** (8003) - Institution management
- **Wallet Service** (8004) - Balance management
- **Payment Service** (8005) - Transaction processing
- **Notification Service** (8006) - Email/SMS alerts
- **Search Service** (8008) - Full-text & semantic search

### Intelligence Layer (Phase 3)
- **AI Service** (8009) - RAG + Claude LLM integration
- **Quiz Service** (8010) - Auto-generating & scoring quizzes
- **Recommendations Service** (8011) - Hybrid content recommendations

### Analytics Layer (Phase 4)
- **Analytics Service** (8012) - Event tracking & reporting

### Infrastructure
- PostgreSQL 17 with pgvector (5432)
- Redis 7 (6379)
- NATS message queue (4222)
- MinIO object storage (9001)
- Typesense search engine (8108)
- Prometheus metrics (9090)
- Grafana dashboards (3000)
- Loki centralized logging (3100)
- Caddy reverse proxy (80/443)

## 📝 Configuration

Edit `.env` before first run:

```bash
# Essential production settings
ENVIRONMENT=production
LOG_LEVEL=info
DB_PASSWORD=your_secure_password
JWT_SECRET=your_long_random_key
CLAUDE_API_KEY=sk-your-actual-key

# Email notifications (optional)
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=app-password

# Domain configuration
CORS_ALLOWED_ORIGINS=https://your-domain.com
```

## 🔌 API Examples

### Create a User
```bash
curl -X POST http://localhost:8000/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "secure_password",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Track an Event (Analytics)
```bash
curl -X POST http://localhost:8012/events/track \
  -H "Content-Type: application/json" \
  -d '{
    "eventType": "view",
    "userId": "user-123",
    "resourceId": "resource-456"
  }'
```

### Generate a Quiz
```bash
curl -X POST http://localhost:8010/quizzes/generate \
  -H "Content-Type: application/json" \
  -d '{
    "resourceId": "resource-456",
    "difficulty": "medium",
    "questionCount": 10
  }'
```

### Ask AI Tutor a Question (Phase 3)
```bash
curl -X POST http://localhost:8009/tutor/ask \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Explain photosynthesis",
    "userId": "user-123",
    "context": "Biology 101"
  }'
```

### Get Recommendations
```bash
curl -X GET 'http://localhost:8011/recommendations/user-123?strategy=hybrid&topK=10'
```

### Get Platform Metrics
```bash
curl http://localhost:8012/metrics/platform?days=30
```

## 📊 Monitoring

| Service | URL | Credentials |
|---------|-----|-------------|
| **Grafana** | http://localhost:3000 | admin / admin |
| **Prometheus** | http://localhost:9090 | - |
| **MinIO** | http://localhost:9001 | minioadmin / minioadmin |
| **Loki** (logs) | http://localhost:3100 | - |

### View Live Metrics
```bash
# All services status
docker compose ps

# Service logs
docker compose logs -f api-gateway

# Database connections
docker compose exec postgres psql -U chengetai -d chengetai -c "\conninfo"

# Redis memory
docker compose exec redis redis-cli info memory
```

## 🛠️ Common Operations

### Health Check All Services
```bash
make health-check
```

### View Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f ai-service

# Last 100 lines
docker compose logs --tail 100 api-gateway
```

### Restart Services
```bash
# Single service
docker compose restart quiz-service

# All services
docker compose restart
```

### Database Operations
```bash
# Connect to database
docker compose exec postgres psql -U chengetai -d chengetai

# Backup database
docker compose exec postgres pg_dump -U chengetai chengetai > backup.sql

# Restore database
docker compose exec postgres psql -U chengetai chengetai < backup.sql

# Check database size
docker compose exec postgres psql -U chengetai -d chengetai \
  -c "SELECT pg_size_pretty(pg_database_size('chengetai'));"
```

### Scale a Service
```bash
# Increase replicas (requires load balancer config)
docker compose up -d --scale quiz-service=3
```

### Stop All Services
```bash
# Keep data
docker compose down

# Remove data (WARNING: irreversible)
docker compose down -v
```

## 🔒 Security Checklist

Before production deployment:

- [ ] Change all default passwords in `.env`
- [ ] Set a strong `JWT_SECRET` (minimum 32 chars)
- [ ] Enable HTTPS (configure Caddy with your domain)
- [ ] Set up firewall rules (allow only necessary ports)
- [ ] Configure database backups
- [ ] Enable monitoring and alerting
- [ ] Set up log retention policies
- [ ] Review and restrict CORS origins
- [ ] Enable rate limiting
- [ ] Configure API authentication for external access

## 📈 Performance Tuning

### Database
```bash
# Increase connection pool
DB_MAX_CONNS=50
DB_MIN_CONNS=10

# Enable query logging
POSTGRES_INITDB_ARGS="-c log_statement=all"
```

### Redis
```bash
# Increase memory limit
command: redis-server --maxmemory 2gb --maxmemory-policy allkeys-lru
```

### Services
```bash
# Adjust log level (info for production)
LOG_LEVEL=info

# Set appropriate batch sizes
ANALYTICS_BATCH_SIZE=100
```

## 🐛 Troubleshooting

### Services Won't Start
```bash
# Check all logs
docker compose logs

# Check specific service
docker compose logs api-gateway

# Verify network
docker network ls | grep chengetai

# Restart everything
docker compose restart
```

### Database Connection Issues
```bash
# Check PostgreSQL status
docker compose ps postgres

# Test connection
docker compose exec postgres psql -U chengetai -d chengetai -c "SELECT 1"

# Check logs
docker compose logs postgres
```

### Port Already in Use
```bash
# Find what's using the port
lsof -i :8000

# Change port in docker-compose.yml or .env
# Then restart: docker compose restart
```

### High Memory Usage
```bash
# Check container stats
docker stats

# Reduce replicas or increase server memory
docker compose down
# Edit docker-compose.yml or increase VM memory
```

## 📚 Full Documentation

For detailed information, see:
- **[DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)** - Complete deployment guide
- **[Makefile](Makefile)** - Available make commands
- **[README.md](README.md)** - Project overview
- **Service READMEs** - Individual service documentation in `apps/*/README.md`

## ✅ Next Steps

1. **Run the deployment**: `./deploy.sh production`
2. **Verify services**: `make health-check`
3. **Access Grafana**: http://localhost:3000
4. **Test APIs**: Use the examples above
5. **Configure domain**: Edit Caddyfile and .env
6. **Set up backups**: Configure PostgreSQL backups
7. **Configure monitoring**: Set up Grafana alerts
8. **Load test**: Verify performance under expected load
9. **User onboarding**: Create initial schools and users
10. **Go live**: DNS + SSL certificate setup

## 🆘 Support

Need help? Check:
1. **Logs**: `docker compose logs -f`
2. **Service README**: `apps/[service]/README.md`
3. **Database schema**: `platform/database/migrations/`
4. **Issues**: https://github.com/wgmasvix-hue/ChengetAi-learn-/issues

---

**Current Status**: ✅ All 13 microservices ready for deployment
**Architecture**: Event-driven, highly scalable, containerized
**Database**: PostgreSQL 17 with pgvector for semantic search
**Monitoring**: Full Prometheus + Grafana stack included
