# ChengetAi Learn - MVP Deployment Guide (Quick Version)

**Minimum Viable Product for Initial Server Launch**

---

## What's Included

✓ PostgreSQL 16 (database)  
✓ Redis 7.2 (cache & sessions)  
✓ Go API Server (Phase One)  
✓ Health checks & monitoring  
✓ Automated backups  
✓ SSL/TLS ready  
✓ Production-grade logging  

**What's NOT included (Phase 2+):**  
- Frontend application (deploy separately)
- AI tutor (Phase 3)
- Curriculum engine (Phase 3)
- Advanced features (later phases)

---

## Deployment Options

### Option A: Automated (Recommended for first-time)

```bash
# 1. SSH into server
ssh user@your-server.com

# 2. Clone repository
git clone https://github.com/wgmasvix-hue/chengetai-learn.git
cd chengetai-learn

# 3. Run automated deployment
sudo bash deploy-mvp.sh

# Script will:
# ✓ Check prerequisites
# ✓ Generate secrets
# ✓ Build Docker images
# ✓ Start services
# ✓ Verify installation
# ✓ Schedule backups
# ✓ Set up monitoring
```

**Time needed:** 5-10 minutes (depending on internet speed)

---

### Option B: Manual (For advanced users)

```bash
# 1. Prepare server
sudo apt update && sudo apt upgrade -y
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# 2. Clone and setup
git clone https://github.com/wgmasvix-hue/chengetai-learn.git
cd chengetai-learn
cp .env.production .env

# 3. Configure secrets (IMPORTANT)
DB_PASSWORD=$(openssl rand -base64 20)
JWT_SECRET=$(openssl rand -base64 32)
sed -i "s/CHANGE_ME_TO_STRONG_PASSWORD_min_20_chars/$DB_PASSWORD/g" .env
sed -i "s/CHANGE_ME_TO_STRONG_SECRET_min_32_chars_random/$JWT_SECRET/g" .env

# 4. Edit domain
nano .env
# Change: CORS_ALLOWED_ORIGINS=https://your-domain.com

# 5. Deploy
docker-compose -f docker-compose.production.yml up -d

# 6. Verify
curl http://localhost:8080/health
```

---

## Architecture

```
┌─────────────┐
│   Nginx     │ (SSL/TLS, reverse proxy)
└──────┬──────┘
       │ :443 → :8080
┌──────▼──────┐
│   Go API    │ (8080)
└──────┬──────┘
       │
   ┌───┴────┐
   │        │
┌──▼──┐ ┌──▼──┐
│ PG  │ │Redis│
│(5432)│ │(6379)│
└─────┘ └─────┘
```

---

## Configuration

### Essential Settings (.env)

```bash
DB_PASSWORD=              # Database password (20+ chars)
JWT_SECRET=               # Authentication secret (32+ chars)
CORS_ALLOWED_ORIGINS=     # Your domain(s)
```

### Generate Secrets (MUST DO)

```bash
# Database password
openssl rand -base64 20

# JWT secret
openssl rand -base64 32

# Copy outputs and paste into .env
```

---

## Testing After Deployment

```bash
# 1. Health check
curl http://localhost:8080/health
# Expected: {"status":"ok"}

# 2. Readiness check (dependencies)
curl http://localhost:8080/ready
# Expected: Shows database and Redis status

# 3. Version info
curl http://localhost:8080/version
```

---

## Common Operations

### View Logs

```bash
# All services
docker-compose -f docker-compose.production.yml logs -f

# API only
docker-compose -f docker-compose.production.yml logs -f api

# Database only
docker-compose -f docker-compose.production.yml logs -f postgres
```

### Stop Services

```bash
docker-compose -f docker-compose.production.yml down
```

### Restart API

```bash
docker-compose -f docker-compose.production.yml restart api
```

### Backup Database

```bash
docker exec chengetai-postgres pg_dump -U chengetai chengetai > backup.sql
```

### Restore Database

```bash
docker exec -i chengetai-postgres psql -U chengetai chengetai < backup.sql
```

---

## Nginx Setup (After API deployment)

```bash
# Install Nginx
sudo apt install -y nginx certbot python3-certbot-nginx

# Create config
sudo nano /etc/nginx/sites-available/chengetai
```

**Paste this:**

```nginx
upstream api {
    server localhost:8080;
}

server {
    listen 80;
    server_name api.your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.your-domain.com;

    ssl_certificate /etc/letsencrypt/live/api.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.your-domain.com/privkey.pem;

    location / {
        proxy_pass http://api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

**Enable and get SSL:**

```bash
sudo ln -s /etc/nginx/sites-available/chengetai /etc/nginx/sites-enabled/
sudo nginx -t
sudo certbot certonly --nginx -d api.your-domain.com
sudo systemctl restart nginx
```

---

## SSL Certificate (Let's Encrypt)

```bash
# Get certificate
sudo certbot certonly --nginx -d api.your-domain.com

# Auto-renew
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

---

## Automatic Startup (Systemd)

```bash
sudo nano /etc/systemd/system/chengetai.service
```

**Paste this:**

```ini
[Unit]
Description=ChengetAi Learn
After=docker.service
Requires=docker.service

[Service]
Type=simple
WorkingDirectory=/home/user/chengetai-learn
ExecStart=/usr/bin/docker-compose -f docker-compose.production.yml up
Restart=always
User=user

[Install]
WantedBy=multi-user.target
```

**Enable:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable chengetai
sudo systemctl start chengetai
```

---

## Monitoring

### Health Check Script

```bash
#!/bin/bash
curl -s http://localhost:8080/health | grep -q "ok" || \
  docker-compose -f docker-compose.production.yml restart api
```

**Add to crontab (every 5 min):**

```bash
crontab -e
# Add: */5 * * * * /home/user/chengetai-learn/monitor.sh
```

### Disk Space

```bash
df -h
# Alert if usage > 80%
```

### Database Size

```bash
docker exec chengetai-postgres du -h /var/lib/postgresql/data
```

---

## Troubleshooting

### API won't start

```bash
# Check logs
docker-compose -f docker-compose.production.yml logs api

# Common causes:
# - Port 8080 in use: lsof -i :8080
# - DB not ready: wait 30 seconds
# - JWT_SECRET too short: must be 32+ chars
```

### Can't connect to database

```bash
# Test connection
docker exec -it chengetai-postgres psql -U chengetai -d chengetai

# Check logs
docker-compose -f docker-compose.production.yml logs postgres
```

### Redis not responding

```bash
# Test Redis
docker exec -it chengetai-redis redis-cli ping
# Expected: PONG
```

### High memory usage

```bash
# Check container stats
docker stats

# Restart services
docker-compose -f docker-compose.production.yml restart
```

---

## Maintenance Schedule

### Daily
- [ ] Check API health: `curl http://localhost:8080/health`
- [ ] Monitor disk space: `df -h`
- [ ] Review logs: `docker-compose logs --tail=100`

### Weekly
- [ ] Verify backups exist
- [ ] Check database size
- [ ] Review error logs

### Monthly
- [ ] Test database restore
- [ ] Update system packages
- [ ] Review and optimize performance
- [ ] Renew SSL cert (automated)

---

## Backup & Recovery

### Automatic Backups
- Scheduled daily at 2 AM
- Location: `/opt/chengetai/backups/`
- Retained: Last 7 days
- Size: ~5-50MB per backup

### Manual Backup

```bash
docker exec chengetai-postgres pg_dump -U chengetai chengetai | gzip > db_backup_$(date +%s).sql.gz
```

### Restore from Backup

```bash
gunzip < db_backup_*.sql.gz | docker exec -i chengetai-postgres psql -U chengetai chengetai
```

---

## Security Checklist

- [ ] JWT_SECRET is strong (32+ chars)
- [ ] DB_PASSWORD is strong (20+ chars)
- [ ] .env file is not in git
- [ ] SSL certificate installed
- [ ] Firewall configured (allow 80, 443 only)
- [ ] SSH key-based auth enabled
- [ ] Password auth disabled
- [ ] Regular backups tested
- [ ] Logs retention configured
- [ ] Health monitoring active

---

## API Endpoints

### Health & Status
```
GET /health           - Liveness probe
GET /ready            - Readiness (checks DB, Redis)
GET /version          - API version info
```

### Authentication (Placeholder)
```
POST /api/v1/auth/register     - User registration
POST /api/v1/auth/login         - User login
GET /api/v1/auth/me             - Current user
```

### Users
```
GET /api/v1/users               - List users
```

---

## Performance Tuning

### For Small Server (1-2 CPU, 2GB RAM)
```bash
# In .env or docker-compose
POSTGRES_INITDB_ARGS="-c max_connections=50"
APP_DB_MAX_CONNS=5
LOG_LEVEL=warn
```

### For Medium Server (4 CPU, 8GB RAM)
```bash
POSTGRES_INITDB_ARGS="-c max_connections=200 -c shared_buffers=512MB"
APP_DB_MAX_CONNS=20
LOG_LEVEL=info
```

### For Large Server (8+ CPU, 32GB+ RAM)
```bash
POSTGRES_INITDB_ARGS="-c max_connections=500 -c shared_buffers=2GB"
APP_DB_MAX_CONNS=50
LOG_LEVEL=info
```

---

## Support & Resources

**Documentation:**
- Full guide: `DEPLOY_MVP.md`
- Architecture: `README.md`
- Security: `SECURITY.md`

**GitHub:**
- Repository: https://github.com/wgmasvix-hue/chengetai-learn
- Issues: Report problems on GitHub

**Logs:**
- API: `docker-compose logs api`
- Database: `docker-compose logs postgres`
- Redis: `docker-compose logs redis`

---

## Production Readiness Checklist

### Infrastructure
- [ ] Server meets minimum requirements (2GB RAM, 10GB disk)
- [ ] Docker and Docker Compose installed
- [ ] Network connectivity verified
- [ ] Ports 80, 443, 8080 available

### Configuration
- [ ] Secrets generated (DB password, JWT secret)
- [ ] Domain configured
- [ ] Email configured (for alerts)
- [ ] Timezone set correctly

### Deployment
- [ ] Services deployed via docker-compose
- [ ] Health checks passing
- [ ] Database initialized
- [ ] Redis running

### Security
- [ ] SSL certificate installed
- [ ] Firewall configured
- [ ] SSH access secured
- [ ] Backups scheduled

### Monitoring
- [ ] Health monitoring active
- [ ] Log aggregation working
- [ ] Alerts configured
- [ ] Backup verification scheduled

### Post-Deployment
- [ ] API endpoints tested
- [ ] Database connected
- [ ] Logging working
- [ ] Backups running

---

## Quick Start Summary

```bash
# 1-command deployment
sudo bash deploy-mvp.sh

# Then verify
curl http://localhost:8080/health

# View logs
docker-compose -f docker-compose.production.yml logs -f
```

**Total time: 5-10 minutes**

---

**Deployment Status:** Ready for Production ✓  
**Version:** MVP (Phase One Foundation)  
**Last Updated:** 2026-08-09

