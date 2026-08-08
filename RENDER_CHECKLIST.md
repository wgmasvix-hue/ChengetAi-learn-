# ChengetAi Platform - Render Deployment Checklist

## Pre-Deployment (Before Creating Services)

### Account Setup
- [ ] Create Render account (https://render.com)
- [ ] Verify email
- [ ] Add payment method (for paid tier)
- [ ] Connect GitHub account
- [ ] Authorize Render to access repositories
- [ ] Set GitHub branch preferences

### Environment Preparation
- [ ] Generate JWT_SECRET: `openssl rand -base64 32`
- [ ] Generate DB_PASSWORD: `openssl rand -base64 16`
- [ ] Gather all API keys:
  - [ ] Claude API key
  - [ ] Stripe API key (if using payments)
  - [ ] SMTP credentials (if using email)
  - [ ] Other third-party keys
- [ ] Prepare domain name (if using custom domain)
- [ ] Document all secrets securely (password manager)

### Repository Preparation
- [ ] Ensure `render.yaml` is in repository root
- [ ] Verify `.env.render` template is present
- [ ] Check Dockerfile is correct: `docker/Dockerfile.service`
- [ ] Verify entrypoint script exists: `docker/entrypoint.sh`
- [ ] Commit and push all changes to GitHub
- [ ] Verify branch is correct (main/master/custom)

## Phase 1: Create Database

### PostgreSQL Setup
- [ ] Go to Render Dashboard
- [ ] Click **New +** → **Database** → **PostgreSQL**
- [ ] Configure:
  - [ ] Database Name: `chengetai-postgres`
  - [ ] Database: `chengetai`
  - [ ] User: `chengetai`
  - [ ] Region: Ohio (or your region)
  - [ ] Plan: Standard ($7/month) or better
- [ ] Click **Create Database**
- [ ] Wait for database to initialize (5-10 minutes)
- [ ] Copy connection details:
  - [ ] Host: `____`
  - [ ] Database: `chengetai`
  - [ ] User: `chengetai`
  - [ ] Password: `____`
  - [ ] External Database URL: `____`

### PostgreSQL Configuration
- [ ] Navigate to database in Render
- [ ] Go to **Info** tab
- [ ] Note Internal Database URL (for docker-compose)
- [ ] Go to **Backups** tab
- [ ] Verify automatic backups enabled
- [ ] Click **Create Backup** to test

## Phase 2: Create Redis Cache

### Redis Setup
- [ ] Go to Render Dashboard
- [ ] Click **New +** → **Cache** → **Redis**
- [ ] Configure:
  - [ ] Cache Name: `chengetai-redis`
  - [ ] Region: Ohio (same as database)
  - [ ] Plan: Standard ($5/month) or better
- [ ] Click **Create Cache**
- [ ] Wait for cache to initialize (2-5 minutes)
- [ ] Copy connection details:
  - [ ] Internal Redis URL: `____`
  - [ ] Connection String: `____`

## Phase 3: Deploy Services Using Blueprint

### Blueprint Deployment
- [ ] Go to Render Dashboard
- [ ] Click **New +** → **Blueprint**
- [ ] Select your GitHub repository
- [ ] Select branch with `render.yaml`
- [ ] Configure:
  - [ ] Name: `chengetai-platform`
  - [ ] Region: Ohio
  - [ ] Environment: Production
- [ ] Add Environment Variables:
  - [ ] JWT_SECRET: `<generated>`
  - [ ] CLAUDE_API_KEY: `<your-key>`
  - [ ] STRIPE_API_KEY: `<your-key>` (if applicable)
  - [ ] CORS_ALLOWED_ORIGINS: `https://your-domain.com`
  - [ ] All other secrets from .env.render
- [ ] Click **Deploy Blueprint**
- [ ] Monitor deployment progress
- [ ] Expected time: 10-20 minutes for all services

### Verify Deployments
- [ ] Check all services show "Live" status
- [ ] Click each service and view logs:
  - [ ] api-gateway logs show "Listening on :8000"
  - [ ] auth-service logs show "Listening on :8001"
  - [ ] (Verify each service)
- [ ] No ERROR messages in logs
- [ ] Services show green health check indicator

## Phase 4: Run Database Migrations

### Migration Execution
- [ ] In Render PostgreSQL dashboard, click **Connect**
- [ ] Copy connection command
- [ ] In your terminal, connect to database:
  ```bash
  psql postgresql://chengetai:<password>@<host>:5432/chengetai
  ```
- [ ] Run migrations in order:
  ```sql
  \i platform/database/migrations/000_create_users_table.sql
  \i platform/database/migrations/001_create_resources_table.sql
  \i platform/database/migrations/002_create_schools_table.sql
  \i platform/database/migrations/003_create_resource_metadata.sql
  \i platform/database/migrations/004_create_user_preferences.sql
  \i platform/database/migrations/005_create_embeddings.sql
  \i platform/database/migrations/006_create_wallet_tables.sql
  \i platform/database/migrations/007_create_quizzes_table.sql
  \i platform/database/migrations/008_create_recommendations_table.sql
  \i platform/database/migrations/009_create_analytics_tables.sql
  ```
- [ ] Verify tables created: `\dt`
- [ ] Exit psql: `\q`

### Verify Schema
- [ ] Check users table exists: `SELECT COUNT(*) FROM users;`
- [ ] Check resources table exists: `SELECT COUNT(*) FROM resources;`
- [ ] Check all tables initialized correctly
- [ ] No errors in output

## Phase 5: Configure Custom Domain

### Domain Setup (Optional but Recommended)
- [ ] Go to API Gateway service in Render
- [ ] Click **Settings** → **Custom Domains**
- [ ] Click **Add Custom Domain**
- [ ] Enter your domain: `api.your-domain.com`
- [ ] Render will provide CNAME record to update
- [ ] Copy CNAME value: `____`

### DNS Configuration
- [ ] Go to your domain registrar (GoDaddy, Namecheap, etc.)
- [ ] Find DNS management section
- [ ] Add CNAME record:
  - [ ] Type: CNAME
  - [ ] Name: `api` (or subdomain)
  - [ ] Value: `<render-provided-cname>`
  - [ ] TTL: 3600
- [ ] Save changes
- [ ] Wait for DNS propagation (up to 24 hours, usually 5-30 mins)
- [ ] Verify in Render dashboard (should show green checkmark)
- [ ] Render auto-provisions SSL certificate (Let's Encrypt)

### SSL Verification
- [ ] Wait for SSL certificate to be issued (2-5 minutes)
- [ ] Test HTTPS: `curl https://api.your-domain.com/health`
- [ ] Should see green SSL indicator in browser

## Phase 6: Health Checks & Testing

### Service Health Checks
```bash
# Test each service
API_URL="https://api.your-domain.com"

# API Gateway
curl $API_URL/health

# Auth Service (using internal routing if available)
curl $API_URL/auth/health

# User Service
curl $API_URL/users/health

# Other services...
```

- [ ] All health endpoints return `{"status":"healthy"}`
- [ ] No connection refused errors
- [ ] No SSL certificate errors

### API Testing
- [ ] Create test user account
- [ ] Test authentication flow
- [ ] Test resource creation
- [ ] Test analytics tracking
- [ ] Test quiz generation (if enabled)
- [ ] Test recommendations

### Load Testing (Optional)
```bash
# Install Apache Bench
# ab -n 100 -c 10 https://api.your-domain.com/health

- [ ] Load test completed
- [ ] Response times acceptable
- [ ] No 5xx errors under load
```

## Phase 7: Production Configuration

### Monitoring & Alerting
- [ ] Set up email alerts:
  - [ ] Service deployment failed
  - [ ] Service crashed
  - [ ] Database error rate high
- [ ] Configure Render dashboard notifications
- [ ] Add monitoring tool (optional):
  - [ ] Datadog
  - [ ] New Relic
  - [ ] Or use Render's built-in metrics

### Logging
- [ ] Verify logs appear in Render dashboard
- [ ] Configure log retention (if available)
- [ ] Set up external logging (optional):
  - [ ] Papertrail
  - [ ] CloudWatch
  - [ ] Datadog Logs

### Backups
- [ ] Enable automatic database backups (should be on)
- [ ] Configure backup retention (7-30 days)
- [ ] Test restore procedure:
  - [ ] Create backup
  - [ ] Restore to new database
  - [ ] Verify data integrity

### Security Hardening
- [ ] All environment variables set (no defaults)
- [ ] No secrets in logs or source code
- [ ] HTTPS enforced (FORCE_HTTPS=true)
- [ ] HSTS enabled (HSTS_ENABLED=true)
- [ ] CORS properly configured (specific domains)
- [ ] Rate limiting enabled
- [ ] Authentication required for protected endpoints

## Phase 8: Go Live

### Pre-Launch
- [ ] All services passing health checks
- [ ] Database fully initialized
- [ ] Custom domain configured and SSL working
- [ ] Backups tested
- [ ] Team access provisioned
- [ ] Documentation updated
- [ ] Support procedures documented

### Launch
- [ ] Update DNS to point to Render custom domain
- [ ] Test from multiple locations
- [ ] Monitor logs for errors
- [ ] Be ready to rollback if needed

### Post-Launch
- [ ] Monitor error rates for 24 hours
- [ ] Check performance metrics
- [ ] Verify backups working
- [ ] Confirm alerts working
- [ ] Document any issues
- [ ] Plan follow-up improvements

## Ongoing Maintenance Checklist

### Daily
- [ ] Check service health status
- [ ] Review error logs
- [ ] Monitor database performance

### Weekly
- [ ] Review analytics
- [ ] Check backup success
- [ ] Review API usage
- [ ] Check for failed deployments

### Monthly
- [ ] Security audit
- [ ] Dependency updates
- [ ] Performance analysis
- [ ] Cost review
- [ ] Disaster recovery test

### Quarterly
- [ ] Full security assessment
- [ ] Capacity planning
- [ ] Architecture review
- [ ] Team training

## Troubleshooting Checklist

### Service Won't Deploy
- [ ] Check render.yaml syntax
- [ ] Verify Dockerfile exists
- [ ] Check GitHub branch is correct
- [ ] Review build logs in Render
- [ ] Verify all environment variables set

### Service Crashes After Deployment
- [ ] View service logs in Render dashboard
- [ ] Check database connectivity
- [ ] Verify environment variables
- [ ] Check for application errors
- [ ] Restart service manually

### Database Connection Error
- [ ] Verify DATABASE_URL format
- [ ] Check database is running (Render dashboard)
- [ ] Verify credentials
- [ ] Test connection manually
- [ ] Check database status page

### High Latency/Performance Issues
- [ ] Check database query performance
- [ ] Review Redis cache hit rate
- [ ] Check service resource usage
- [ ] Scale up instance if needed
- [ ] Enable CDN if available

### SSL Certificate Issues
- [ ] Verify DNS CNAME is correct
- [ ] Wait for DNS propagation
- [ ] Check domain ownership
- [ ] Review Render certificate status
- [ ] Request manual certificate renewal

## Cost Optimization Checklist

- [ ] Choose appropriate plan size
- [ ] Scale down unused services
- [ ] Enable auto-scaling (if available)
- [ ] Monitor resource usage
- [ ] Consider reserved instances
- [ ] Review and cancel unused services

## Security Checklist

- [ ] All secrets in Render environment variables
- [ ] No secrets in code or logs
- [ ] HTTPS enforced
- [ ] Database SSL required
- [ ] Backups encrypted
- [ ] Access control configured
- [ ] Audit logging enabled
- [ ] Rate limiting active
- [ ] CORS restrictive
- [ ] Secrets rotated regularly

---

**Deployment Status**: Ready for Render ✅  
**Last Updated**: 2026-08-06  
**Next Review**: Post-deployment + 1 week
