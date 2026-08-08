# ChengetAi Platform - Render Deployment Guide

## Overview

This guide covers deploying the ChengetAi Platform on Render, a modern cloud platform for hosting applications.

## Architecture on Render

```
Render Services:
├── PostgreSQL (Managed Database)
├── Redis (Managed Cache)
├── API Gateway (Web Service) - Entry point
├── Auth Service (Web Service)
├── User Service (Web Service)
├── School Service (Web Service)
├── Wallet Service (Web Service)
├── Payment Service (Web Service)
├── Notification Service (Web Service)
├── Search Service (Web Service)
├── AI Service (Web Service) - Phase 3
├── Quiz Service (Web Service) - Phase 3
├── Recommendations Service (Web Service) - Phase 3
└── Analytics Service (Web Service) - Phase 4
```

## Prerequisites

1. **Render Account**: https://render.com (free tier available)
2. **GitHub Repository**: Connected to Render
3. **Environment Variables**: Prepared in advance
4. **Credit Card**: For paid services (free tier limited to 3 services + 1 database)

## Step 1: Prepare Environment Variables

Before deploying, gather all required secrets:

```bash
# Generate strong secrets
JWT_SECRET=$(openssl rand -base64 32)
DB_PASSWORD=$(openssl rand -base64 16)
REDIS_PASSWORD=$(openssl rand -base64 16)

echo "JWT_SECRET=$JWT_SECRET"
echo "DB_PASSWORD=$DB_PASSWORD"
echo "REDIS_PASSWORD=$REDIS_PASSWORD"
```

**Critical Environment Variables** (add to Render):
```
ENVIRONMENT=production
LOG_LEVEL=info
JWT_SECRET=<generated_value>
DB_PASSWORD=<generated_value>
REDIS_PASSWORD=<generated_value>
CORS_ALLOWED_ORIGINS=https://your-domain.com
CLAUDE_API_KEY=sk-your-key
STRIPE_API_KEY=<stripe_key>
```

## Step 2: Create PostgreSQL Database

1. Go to **Dashboard** → **Databases** → **Create Database**
2. Select **PostgreSQL**
3. Configure:
   - Name: `chengetai-postgres`
   - Database: `chengetai`
   - User: `chengetai`
   - Region: `Ohio` (or your preference)
   - Plan: Standard ($7/month)
4. Click **Create Database**
5. Note the connection string and credentials

## Step 3: Create Redis Cache

1. Go to **Dashboard** → **Caches** → **Create Cache**
2. Select **Redis**
3. Configure:
   - Name: `chengetai-redis`
   - Region: `Ohio`
   - Plan: Standard ($5/month)
4. Click **Create Cache**
5. Note the connection string

## Step 4: Deploy Services

### Option A: Using render.yaml (Recommended)

Render supports infrastructure-as-code via `render.yaml` included in repository.

1. Go to **Dashboard** → **New +** → **Blueprint**
2. Select your GitHub repository
3. Branch: `main` (or your branch)
4. Region: `Ohio`
5. Name: `chengetai-platform`
6. Click **Deploy Blueprint**

Render will:
- Parse `render.yaml`
- Create all services
- Set up environment variables
- Configure databases
- Deploy all microservices

**Note**: Free tier limits 3 services + 1 database. Upgrade to paid plan for all 12 services.

### Option B: Manual Service Creation

If `render.yaml` doesn't work, create services manually:

1. **API Gateway** (Entry point)
   - Service Type: Web Service
   - Name: `chengetai-api-gateway`
   - Runtime: Docker
   - Repository: Your GitHub repo
   - Docker: `docker/Dockerfile.service`
   - Build Command: `docker build --build-arg SERVICE_NAME=api-gateway -t app .`
   - Start Command: `./app`
   - Port: 8000
   - Environment Variables:
     - SERVICE_NAME: `api-gateway`
     - SERVICE_PORT: `8000`
     - DB_HOST: (PostgreSQL host)
     - DB_USER: `chengetai`
     - DB_PASSWORD: (from database)
     - DB_NAME: `chengetai`

2. Repeat for each service (auth-service, user-service, etc.)
   - Change SERVICE_NAME and SERVICE_PORT for each
   - All share same PostgreSQL database and Redis cache

## Step 5: Configure Custom Domain

1. Go to Service → **Settings** → **Custom Domain**
2. Add your domain: `api.your-domain.com`
3. Update DNS records:
   ```
   Type: CNAME
   Name: api
   Value: <render_provided_cname>
   ```
4. Render auto-provisions SSL certificate (Let's Encrypt)

## Step 6: Run Database Migrations

After database is created:

```bash
# Connect to PostgreSQL
psql postgresql://chengetai:<password>@<host>:5432/chengetai

# Run migrations
\i platform/database/migrations/000_create_users_table.sql
\i platform/database/migrations/001_create_resources_table.sql
... (run all SQL files in order)
```

Or use Render's PostgreSQL dashboard SQL editor.

## Step 7: Verify Deployment

Check each service:

```bash
# API Gateway health check
curl https://api.your-domain.com/health

# Each service
curl https://api.your-domain.com:8001/health  # Auth Service
curl https://api.your-domain.com:8002/health  # User Service
... (check each port)
```

## Pricing on Render

| Service | Type | Cost | Notes |
|---------|------|------|-------|
| PostgreSQL | Database | $7/month | Managed, automated backups |
| Redis | Cache | $5/month | Managed |
| Web Service (each) | Compute | $7/month | OR use free tier (limited) |
| **Total (12 services)** | | **~$100/month** | Scales with usage |

**Free Tier**:
- 1 PostgreSQL instance (limited)
- Up to 3 web services
- 100 GB/month egress
- Perfect for testing, limited for production

**Paid Tier** (Recommended):
- Unlimited services
- Professional support
- Advanced monitoring
- Auto-scaling

## Environment-Specific Configuration

### Production (render.yaml)
- All services deployed
- Auto SSL certificates
- Database backups enabled
- Logging to Render dashboard
- Health checks enabled

### Staging
- Subset of services
- Same database structure
- Different credentials
- Lower resource allocation

### Development
- Single service (API Gateway)
- Separate development database
- Mock services for testing

## Monitoring & Logs

1. **Service Logs**: Dashboard → Service → Logs
2. **Metrics**: Dashboard → Service → Metrics
3. **Database**: Dashboard → Database → Metrics
4. **Alerts**: Configure email alerts for failures

## Deployment Workflow

```bash
# 1. Push changes to GitHub
git add .
git commit -m "Deploy to Render"
git push origin main

# 2. Render auto-deploys on push (if auto-deploy enabled)
# 3. Check deployment status in Render dashboard
# 4. Monitor logs for errors
# 5. Test endpoints

# Rolling back
git revert <commit>
git push origin main
# Render auto-redeploys previous version
```

## Troubleshooting

### Service Won't Start
1. Check logs: Dashboard → Logs
2. Verify environment variables set
3. Check port matches SERVICE_PORT
4. Verify database connectivity

### Database Connection Error
1. Check DATABASE_URL format
2. Verify firewall rules (Render allows all)
3. Test connection:
   ```bash
   psql $DATABASE_URL
   ```

### Health Check Failures
1. Verify /health endpoint exists
2. Check port is correct
3. Review application logs

### Performance Issues
1. Scale up instance: Settings → Plan
2. Increase database resources
3. Add Redis caching
4. Enable CDN for static assets

## Disaster Recovery

### Database Backup
- Render automatically backs up PostgreSQL daily
- Retained for 7 days
- Manual backups: Database → Backups

### Restore Database
1. Create new PostgreSQL instance
2. Restore from backup
3. Update connection strings
4. Restart services

## Cost Optimization

1. **Use free tier for development**
2. **Combine services** into single container if possible
3. **Scale down** unused services
4. **Use spot instances** (Render Spot, when available)
5. **Monitor resource usage** and right-size

## Security Checklist

- [ ] All environment variables encrypted
- [ ] Database SSL enforced
- [ ] CORS configured correctly
- [ ] Rate limiting enabled
- [ ] Authentication required for endpoints
- [ ] Secrets not in code or logs
- [ ] HTTPS enforced
- [ ] Firewall rules configured
- [ ] Backups enabled and tested
- [ ] Monitoring/alerting configured

## Going Live Checklist

- [ ] Domain configured and DNS updated
- [ ] SSL certificate verified
- [ ] Database migrations completed
- [ ] All services deployed and healthy
- [ ] Environment variables correct
- [ ] Backups working
- [ ] Monitoring alerts configured
- [ ] Team access provisioned
- [ ] Documentation updated
- [ ] Load testing completed

## Support & Resources

- **Render Docs**: https://render.com/docs
- **Render Support**: https://support.render.com
- **Platform README**: See README.md in repository
- **Security Guide**: See SECURITY.md

## Next Steps

1. Create Render account
2. Connect GitHub repository
3. Deploy using render.yaml
4. Configure custom domain
5. Run database migrations
6. Test all endpoints
7. Set up monitoring
8. Document access credentials (secure storage)

---

**Deployment Status**: Ready for Render ✅  
**Last Updated**: 2026-08-06  
**Maintainer**: ChengetAi Team
