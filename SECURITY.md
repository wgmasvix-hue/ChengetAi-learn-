# ChengetAi Platform - Security Documentation

## Overview

This document outlines security practices, configuration, and hardening guidelines for the ChengetAi Platform.

## Table of Contents

1. [Authentication & Authorization](#authentication--authorization)
2. [Secrets Management](#secrets-management)
3. [Database Security](#database-security)
4. [API Security](#api-security)
5. [Infrastructure Security](#infrastructure-security)
6. [Development Security](#development-security)
7. [Incident Response](#incident-response)
8. [Security Checklist](#security-checklist)

---

## Authentication & Authorization

### JWT Implementation

The platform uses JWT (JSON Web Tokens) for stateless authentication:

```go
// Generate tokens with controlled expiry
accessToken, _ := jwtManager.GenerateAccessToken(userID, email, role, schoolID)
// Access tokens: 15 minutes (short-lived)
// Refresh tokens: 7 days (long-lived)
```

**Security Properties:**
- ✅ HMAC-SHA256 signing algorithm
- ✅ Issuer and audience validation
- ✅ Expiry time enforcement
- ✅ Token revocation support via Redis blacklist

**Best Practices:**
```go
// Always validate against current issuer/audience
if !claims.VerifyAudience("chengetai-platform", true) {
    return nil, fmt.Errorf("invalid token audience")
}

// Check signing method to prevent algorithm confusion attacks
if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
    return nil, fmt.Errorf("unexpected signing method")
}
```

### Password Security

Passwords are hashed using bcrypt with configurable cost:

```go
// Minimum requirements
- Length: 8-72 characters
- Upper + lowercase + numbers required
- Bcrypt cost: 12 (configurable via BCRYPT_COST)

// Validate before hashing
if err := IsValidPassword(password); err != nil {
    return err // "password must contain uppercase, lowercase, and numbers"
}

// Hash with bcrypt
hasher := NewPasswordHasher(bcrypt.DefaultCost)
hash, err := hasher.Hash(password)
```

**Configuration:**
```env
PASSWORD_MIN_LENGTH=8
PASSWORD_REQUIRE_UPPERCASE=true
PASSWORD_REQUIRE_LOWERCASE=true
PASSWORD_REQUIRE_DIGITS=true
BCRYPT_COST=12
```

---

## Secrets Management

### Environment Variables

**CRITICAL: Never commit secrets to version control**

```bash
# Generate strong secrets
openssl rand -base64 32  # 32-byte base64 for JWT
openssl rand -base64 48  # 48-byte base64 for passwords
```

**Required Secrets (minimum 32 characters):**
```env
JWT_SECRET=<generated_random_32_chars>
DB_PASSWORD=<strong_password>
REDIS_PASSWORD=<strong_password>
MINIO_ROOT_PASSWORD=<strong_password>
TYPESENSE_API_KEY=<strong_api_key>
```

### Runtime Validation

All required secrets are validated at startup:

```go
secrets := NewSecretsManager()

// Validates all critical secrets exist and meet minimum requirements
if err := secrets.ValidateSecrets(); err != nil {
    panic(err) // Fails fast on misconfiguration
}

// Additional production validation
if err := secrets.VerifyProductionSecrets(); err != nil {
    panic(err) // Production-only checks
}
```

**Production Verification:**
```
✓ JWT_SECRET >= 32 characters
✓ DB_PASSWORD >= 16 characters
✓ REDIS_PASSWORD >= 16 characters
✓ CORS_ALLOWED_ORIGINS not set to "*"
✓ ENVIRONMENT set to "production"
```

---

## Database Security

### Connection Security

**Development (SSL disabled):**
```go
dsn := "postgres://user:pass@localhost:5432/db?sslmode=disable"
// ONLY for development
```

**Production (SSL required):**
```go
dsn := "postgres://user:pass@host:5432/db?sslmode=require"
// Forces encrypted connections
```

### SQL Injection Prevention

All database queries use parameterized statements:

```go
// ✅ CORRECT: Parameterized query
db.QueryRow("SELECT * FROM users WHERE id = $1", userID).Scan(...)

// ❌ WRONG: String interpolation
db.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID))
```

**Input Validation:**
```go
validator := NewInputValidator()

// Validate identifiers before using in queries
if err := validator.ValidateSQLIdentifier(tableName); err != nil {
    return err // Prevents SQL injection
}

// Check for obvious SQL injection patterns
if err := validator.CheckSQLInjectionPattern(userInput); err != nil {
    return err
}
```

### Connection Pooling

```env
DB_MAX_CONNS=25      # Maximum connections
DB_MIN_CONNS=5       # Minimum idle connections
DB_POOL_SIZE=10      # Initial pool size
```

---

## API Security

### Security Headers Middleware

All HTTP responses include security headers:

```
X-Content-Type-Options: nosniff              # Prevent MIME sniffing
X-Frame-Options: DENY                        # Prevent clickjacking
X-XSS-Protection: 1; mode=block             # XSS protection
Strict-Transport-Security: max-age=31536000 # Force HTTPS
Content-Security-Policy: default-src 'self' # Restrict resources
```

### CORS Configuration

**Development:**
```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

**Production (NEVER wildcard):**
```env
CORS_ALLOWED_ORIGINS=https://chengetai.com,https://www.chengetai.com
```

**Validation:**
```go
// In production, CORS_ALLOWED_ORIGINS cannot be "*"
if environment == "production" && corsOrigins == "*" {
    panic("CORS wildcard not allowed in production")
}
```

### Rate Limiting

```env
RATE_LIMIT_REQUESTS=100   # Requests per window
RATE_LIMIT_WINDOW=1m      # Time window
```

Implemented at middleware level:
```go
// Apply to all endpoints
rateLimiter := NewRateLimiter(100, 1*time.Minute)
router.Use(RateLimitMiddlewareFunc(rateLimiter))
```

### Input Validation

All user inputs validated on entry:

```go
validator := NewInputValidator()

// Email validation
if err := validator.ValidateEmail(email); err != nil {
    return http.StatusBadRequest, err
}

// Username validation (3-32 chars, alphanumeric + hyphen/underscore)
if err := validator.ValidateUsername(username); err != nil {
    return http.StatusBadRequest, err
}

// Generic string validation with length limits
if err := validator.ValidateString("bio", bio, 0, 500); err != nil {
    return http.StatusBadRequest, err
}

// XSS prevention via sanitization
cleanedInput := validator.SanitizeString(userInput)
```

---

## Infrastructure Security

### Docker Security

**Dockerfile Best Practices:**
```dockerfile
# Use specific version tags (not 'latest')
FROM golang:1.21-alpine3.18

# Run as non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Use multi-stage builds to reduce image size
FROM golang:1.21 as builder
# ... build steps

FROM alpine:3.18
COPY --from=builder /app/bin /app/
```

### Network Security

**Docker Compose Network:**
```yaml
networks:
  chengetai:
    driver: bridge
    # Services only accessible within network
```

**Reverse Proxy (Caddy):**
```
- Handles HTTPS termination
- Rate limiting at edge
- Request validation
- Request/response logging
```

### Database Backups

```bash
# Automated daily backups
docker compose exec postgres pg_dump -U chengetai chengetai > backup.sql

# Encrypted backup storage
gpg --encrypt backup.sql

# Verify backup integrity
pg_restore --list backup.sql
```

---

## Development Security

### Code Review Checklist

- [ ] No hardcoded secrets or credentials
- [ ] Parameterized database queries used
- [ ] Input validation on all endpoints
- [ ] Error messages don't leak sensitive info
- [ ] No debugging/profiling endpoints in production
- [ ] Dependencies are from trusted sources
- [ ] CORS properly configured for environment
- [ ] Authentication required for sensitive endpoints
- [ ] Rate limiting applied to public endpoints
- [ ] Logging redacts sensitive data

### Dependency Management

```bash
# Audit Go dependencies for vulnerabilities
go list -json -m all | nancy sleuth

# Update to latest secure versions
go get -u ./...
go mod tidy

# Check for known CVEs
go run github.com/aquasecurity/trivy@latest fs ./
```

### Git Security

**Prevent accidental secret commits:**
```bash
# Install pre-commit hooks
pre-commit install

# Create .pre-commit-config.yaml
repos:
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.4.0
    hooks:
      - id: detect-secrets
        args: ['--baseline', '.secrets.baseline']
```

---

## Incident Response

### Potential Security Events

1. **Compromised API Key**
   - Revoke immediately in provider console
   - Rotate JWT secret
   - Review access logs
   - Notify users if data accessed

2. **Database Breach**
   - Enable audit logging
   - Review and export sensitive data access logs
   - Reset all passwords
   - Enable 2FA for admin accounts

3. **DDoS Attack**
   - Increase rate limits temporarily
   - Block attacking IP ranges at firewall
   - Scale infrastructure
   - Contact DDoS mitigation service

### Security Logging

All security events logged:
```
- Failed authentication attempts
- Rate limit violations
- Unusual database access patterns
- Configuration changes
- Secret access/rotation
```

Logs sent to Loki with PII redaction:
```env
ANALYTICS_REDACT_PII=true  # Remove email, phone, etc from logs
LOG_LEVEL=info             # Production: info, Development: debug
```

---

## Security Checklist

### Pre-Deployment

- [ ] All environment secrets set to strong values (32+ chars)
- [ ] `ENVIRONMENT` set to `production`
- [ ] `DB_SSL_MODE` set to `require`
- [ ] `FORCE_HTTPS` enabled
- [ ] `CORS_ALLOWED_ORIGINS` set to specific domains
- [ ] `JWT_SECRET` generated with `openssl rand -base64 32`
- [ ] Database backups configured and tested
- [ ] SMTP credentials configured for alerts
- [ ] Monitoring/logging configured (Prometheus, Loki)
- [ ] Rate limiting configured appropriately
- [ ] Admin password changed from default
- [ ] Grafana password changed from default
- [ ] MinIO credentials changed from default
- [ ] SSL certificates generated/updated
- [ ] Firewall rules configured (allow only necessary ports)
- [ ] Secrets manager or vault configured

### Ongoing Maintenance

- [ ] Weekly: Review authentication logs for suspicious activity
- [ ] Monthly: Run dependency vulnerability scan (`nancy`, `trivy`)
- [ ] Monthly: Review database access patterns
- [ ] Quarterly: Security audit of critical components
- [ ] Quarterly: Update all dependencies
- [ ] Semi-annually: Penetration testing
- [ ] Annually: Full security assessment

### Post-Deployment

- [ ] Monitor error logs for security-related errors
- [ ] Set up alerts for failed authentication attempts
- [ ] Set up alerts for rate limit violations
- [ ] Set up alerts for database errors
- [ ] Verify HTTPS is working correctly
- [ ] Test password reset flow
- [ ] Test 2FA implementation (when deployed)
- [ ] Verify backups are working

---

## Security Best Practices

### For Developers

1. **Never commit secrets** - Use `.env` files, ignore from git
2. **Validate all inputs** - Sanitize and validate on entry
3. **Use parameterized queries** - Never interpolate values
4. **Log security events** - Track failed auth, rate limits, etc
5. **Principle of least privilege** - Users only get needed permissions
6. **Defense in depth** - Multiple layers of security controls
7. **Fail securely** - Errors don't expose sensitive information
8. **Keep dependencies updated** - Regular security patches

### For DevOps

1. **Rotate secrets regularly** - Every 90 days minimum
2. **Monitor access patterns** - Look for anomalies
3. **Encrypt in transit** - HTTPS/TLS everywhere
4. **Encrypt at rest** - For sensitive data in database
5. **Backup regularly** - Test restore procedures
6. **Network segmentation** - Limit service-to-service communication
7. **Audit everything** - Track all system changes
8. **Incident response plan** - Know your procedures

---

## Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://pkg.go.dev/std#security)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/sql-syntax.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc7519)
- [Docker Security](https://docs.docker.com/engine/security/)

---

**Last Updated**: 2026-08-06  
**Security Contact**: security@chengetai.africa
