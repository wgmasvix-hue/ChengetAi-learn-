# ChengetAi Learn - MVP Production Deployment

**This guide deploys the minimum viable version for initial server launch.**

Includes only essentials:
- PostgreSQL database
- Redis cache
- Go API server
- Basic authentication
- Health checks

**Estimated deployment time: 5-10 minutes**

---

## Prerequisites

- Ubuntu 20.04 LTS or later
- Docker & Docker Compose installed
- 2GB+ RAM, 10GB+ disk space
- Open ports: 80, 443, 8080 (API)

### Install Docker

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
```

### Verify Installation

```bash
docker --version
docker-compose --version
```

---

## Step 1: Prepare Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Create app directory
mkdir -p /opt/chengetai
cd /opt/chengetai

# Clone repository (or download files)
git clone https://github.com/wgmasvix-hue/chengetai-learn.git .
# OR
# scp -r * user@server:/opt/chengetai/
```

---

## Step 2: Configure Environment

```bash
# Copy production environment
cp .env.production .env

# Generate strong secrets
DB_PASSWORD=$(openssl rand -base64 20)
JWT_SECRET=$(openssl rand -base64 32)

# Edit .env with your values
nano .env
```

**Required changes in .env:**

```bash
DB_PASSWORD=<generated-20-char-password>
JWT_SECRET=<generated-32-char-secret>
CORS_ALLOWED_ORIGINS=https://your-domain.com
```

---

## Step 3: Build & Launch

```bash
# Build Docker image
docker-compose -f docker-compose.production.yml build

# Start services (runs in background)
docker-compose -f docker-compose.production.yml up -d

# Check status
docker-compose -f docker-compose.production.yml ps
```

**Expected output:**
```
NAME                 STATUS              PORTS
chengetai-postgres   Up (healthy)        5432/tcp
chengetai-redis      Up (healthy)        6379/tcp
chengetai-api        Up (healthy)        0.0.0.0:8080->8080/tcp
```

---

## Step 4: Verify Installation

```bash
# Check API health
curl http://localhost:8080/health
# Expected: {"status":"ok"}

# Check readiness (dependencies)
curl http://localhost:8080/ready
# Expected: Shows database and redis status

# Check version
curl http://localhost:8080/version
```

---

## Step 5: Configure Web Server (Nginx)

```bash
# Install Nginx
sudo apt install -y nginx

# Create config
sudo nano /etc/nginx/sites-available/chengetai
```

**Add this config:**

```nginx
upstream chengetai_api {
    server localhost:8080;
}

server {
    listen 80;
    server_name api.your-domain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.your-domain.com;

    ssl_certificate /etc/letsencrypt/live/api.your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.your-domain.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://chengetai_api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        proxy_pass http://chengetai_api;
        access_log off;
    }
}
```

**Enable site:**

```bash
sudo ln -s /etc/nginx/sites-available/chengetai /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

---

## Step 6: SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Get certificate
sudo certbot certonly --nginx -d api.your-domain.com

# Auto-renewal
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

---

## Step 7: Database Initialization

Migrations run automatically on startup. Verify:

```bash
# Connect to database
docker exec -it chengetai-postgres psql -U chengetai -d chengetai

# List tables
\dt

# Exit
\q
```

**Expected tables:**
- users
- learner_profiles
- teacher_profiles

---

## Step 8: Monitoring & Logs

```bash
# View logs
docker-compose -f docker-compose.production.yml logs -f

# API logs only
docker-compose -f docker-compose.production.yml logs -f api

# Database logs
docker-compose -f docker-compose.production.yml logs -f postgres

# Redis logs
docker-compose -f docker-compose.production.yml logs -f redis
```

---

## Step 9: Backup Strategy

```bash
# Create backup script
cat > /opt/chengetai/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/chengetai/backups"
mkdir -p $BACKUP_DIR
DATE=$(date +%Y%m%d_%H%M%S)

# Backup PostgreSQL
docker exec chengetai-postgres pg_dump -U chengetai chengetai > $BACKUP_DIR/db_$DATE.sql

# Compress
gzip $BACKUP_DIR/db_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -type f -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/db_$DATE.sql.gz"
EOF

chmod +x /opt/chengetai/backup.sh

# Schedule daily backups at 2 AM
crontab -e
# Add: 0 2 * * * /opt/chengetai/backup.sh
```

---

## Step 10: Automated Restart

```bash
# Create systemd service
sudo nano /etc/systemd/system/chengetai.service
```

**Add this:**

```ini
[Unit]
Description=ChengetAi Learn
After=docker.service
Requires=docker.service

[Service]
Type=simple
WorkingDirectory=/opt/chengetai
ExecStart=/usr/bin/docker-compose -f docker-compose.production.yml up
ExecStop=/usr/bin/docker-compose -f docker-compose.production.yml down
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable service:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable chengetai
sudo systemctl start chengetai
```

---

## Monitoring

### Health Check

```bash
# Create monitoring script
cat > /opt/chengetai/monitor.sh << 'EOF'
#!/bin/bash
STATUS=$(curl -s http://localhost:8080/health | grep "ok" | wc -l)
if [ $STATUS -eq 0 ]; then
  echo "API DOWN - $(date)" | mail -s "ChengetAi Alert" admin@domain.com
  docker-compose -f docker-compose.production.yml restart api
fi
EOF

# Run every 5 minutes
crontab -e
# Add: */5 * * * * /opt/chengetai/monitor.sh
```

### Disk Space

```bash
# Check disk usage
df -h

# Check Docker volumes
docker volume ls
docker volume inspect chengetai_postgres_data
```

---

## Common Commands

```bash
# Start services
docker-compose -f docker-compose.production.yml up -d

# Stop services
docker-compose -f docker-compose.production.yml down

# Restart API
docker-compose -f docker-compose.production.yml restart api

# View logs
docker-compose -f docker-compose.production.yml logs -f api

# Connect to database
docker exec -it chengetai-postgres psql -U chengetai -d chengetai

# Backup database
docker exec chengetai-postgres pg_dump -U chengetai chengetai > backup.sql

# Update API (rebuild)
docker-compose -f docker-compose.production.yml build --no-cache
docker-compose -f docker-compose.production.yml up -d
```

---

## Troubleshooting

### API won't start

```bash
# Check logs
docker-compose -f docker-compose.production.yml logs api

# Common issues:
# - Port 8080 already in use: change APP_PORT in .env
# - Database not ready: wait 30s and check `docker ps`
# - JWT_SECRET too short: must be 32+ characters
```

### Database connection error

```bash
# Test connection
docker exec -it chengetai-postgres psql -U chengetai -d chengetai -c "SELECT 1"

# Check database logs
docker-compose -f docker-compose.production.yml logs postgres
```

### Redis connection error

```bash
# Test Redis
docker exec -it chengetai-redis redis-cli ping
# Expected: PONG

# Check Redis logs
docker-compose -f docker-compose.production.yml logs redis
```

### High disk usage

```bash
# Prune Docker system
docker system prune -a --volumes -f

# Check log sizes
docker-compose -f docker-compose.production.yml logs --tail=0

# Limit log size in compose file (already configured)
```

---

## Upgrade Strategy

```bash
# 1. Pull latest code
cd /opt/chengetai
git pull origin main

# 2. Rebuild image
docker-compose -f docker-compose.production.yml build

# 3. Restart services (zero downtime)
docker-compose -f docker-compose.production.yml up -d

# 4. Verify
curl https://api.your-domain.com/health
```

---

## Security Checklist

- [ ] JWT_SECRET is strong (32+ random characters)
- [ ] DB_PASSWORD is strong (20+ random characters)
- [ ] CORS_ALLOWED_ORIGINS set to actual domain (not *)
- [ ] SSL certificate installed and auto-renewing
- [ ] Backups running daily
- [ ] Firewall configured (only 80, 443 open)
- [ ] Database backups tested (restore one)
- [ ] Log retention configured
- [ ] Health checks running

---

## Performance Tuning

**PostgreSQL** (in .env or docker-compose):
```yaml
environment:
  POSTGRES_INITDB_ARGS: "-c max_connections=200 -c shared_buffers=256MB"
```

**Redis** (already optimized for production)

**API** (adjust if needed):
```bash
APP_PORT=8080          # Change if port conflict
DB_MAX_CONNS=10        # Increase for high traffic
DB_MIN_CONNS=2
LOG_LEVEL=info         # Change to 'warn' for less logging
```

---

## Deployment Complete ✓

**API URL:** https://api.your-domain.com
**Health Check:** https://api.your-domain.com/health
**Status:** https://api.your-domain.com/ready

---

## Next Steps

1. **Add Frontend** - Deploy React/Next.js application
2. **Enable Payments** - Configure Stripe if needed
3. **User Onboarding** - Create first admin account
4. **Monitoring** - Set up alerts and dashboards
5. **Scale** - Add more replicas as traffic grows

---

**Deployment Date:** _____________  
**Deployed By:** _____________  
**Domain:** _____________  
**Status:** ✓ Production Ready

