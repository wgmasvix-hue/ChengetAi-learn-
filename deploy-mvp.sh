#!/bin/bash
# ChengetAi Learn - MVP Deployment Script
# Automated setup for production server deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}========================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"

    # Check if running as root
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root or with sudo"
        exit 1
    fi
    print_success "Running as root"

    # Check Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker not installed"
        exit 1
    fi
    print_success "Docker installed: $(docker --version)"

    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose not installed"
        exit 1
    fi
    print_success "Docker Compose installed: $(docker-compose --version)"

    # Check disk space
    AVAILABLE=$(df / | awk 'NR==2 {print $4}')
    if [ "$AVAILABLE" -lt 10485760 ]; then  # 10GB in KB
        print_warning "Less than 10GB available disk space"
    else
        print_success "Disk space available: $(numfmt --to=iec $AVAILABLE 2>/dev/null || echo "$AVAILABLE KB")"
    fi

    # Check ports
    for port in 80 443 5432 6379 8080; do
        if netstat -tuln 2>/dev/null | grep -q ":$port "; then
            print_warning "Port $port already in use"
        fi
    done
}

# Setup directory structure
setup_directories() {
    print_header "Setting Up Directories"

    mkdir -p /opt/chengetai/{backups,logs,configs}
    print_success "Created directories"

    # Create Docker volumes directory
    mkdir -p /var/lib/docker/volumes
    print_success "Docker volumes prepared"
}

# Generate secrets
generate_secrets() {
    print_header "Generating Secrets"

    DB_PASSWORD=$(openssl rand -base64 20)
    JWT_SECRET=$(openssl rand -base64 32)

    print_success "Generated database password"
    print_success "Generated JWT secret"

    echo "DB_PASSWORD=$DB_PASSWORD" > /tmp/secrets.txt
    echo "JWT_SECRET=$JWT_SECRET" >> /tmp/secrets.txt
}

# Configure environment
configure_environment() {
    print_header "Configuring Environment"

    if [ -f .env.production ]; then
        cp .env.production .env
        print_success "Copied .env.production to .env"
    else
        print_error ".env.production not found"
        exit 1
    fi

    # Replace secrets in .env
    sed -i "s/CHANGE_ME_TO_STRONG_PASSWORD_min_20_chars/$DB_PASSWORD/g" .env
    sed -i "s/CHANGE_ME_TO_STRONG_SECRET_min_32_chars_random/$JWT_SECRET/g" .env
    print_success "Configured secrets"

    # Ask for domain
    read -p "Enter your domain (e.g., api.yourdomain.com): " DOMAIN
    sed -i "s|https://your-domain.com,https://api.your-domain.com|https://$DOMAIN|g" .env
    print_success "Configured domain: $DOMAIN"

    # Ask for app environment
    read -p "Enter APP_ENV (production/staging): " APP_ENV
    sed -i "s/production/$APP_ENV/g" .env
    print_success "Configured APP_ENV: $APP_ENV"
}

# Build Docker images
build_images() {
    print_header "Building Docker Images"

    docker-compose -f docker-compose.production.yml build --no-cache
    print_success "Docker images built"
}

# Start services
start_services() {
    print_header "Starting Services"

    docker-compose -f docker-compose.production.yml up -d
    print_success "Services started"

    # Wait for services to be ready
    echo "Waiting for services to be healthy..."
    sleep 10

    # Check status
    docker-compose -f docker-compose.production.yml ps
}

# Verify installation
verify_installation() {
    print_header "Verifying Installation"

    # Test PostgreSQL
    if docker exec chengetai-postgres pg_isready -U chengetai &> /dev/null; then
        print_success "PostgreSQL is running"
    else
        print_error "PostgreSQL failed"
        exit 1
    fi

    # Test Redis
    if docker exec chengetai-redis redis-cli ping &> /dev/null; then
        print_success "Redis is running"
    else
        print_error "Redis failed"
        exit 1
    fi

    # Test API
    sleep 5
    if curl -s http://localhost:8080/health | grep -q "ok"; then
        print_success "API is responding"
    else
        print_warning "API not responding yet, may take a moment..."
    fi
}

# Setup backup script
setup_backups() {
    print_header "Setting Up Backups"

    cat > /opt/chengetai/backup.sh << 'BACKUP_EOF'
#!/bin/bash
BACKUP_DIR="/opt/chengetai/backups"
mkdir -p $BACKUP_DIR
DATE=$(date +%Y%m%d_%H%M%S)

docker exec chengetai-postgres pg_dump -U chengetai chengetai > $BACKUP_DIR/db_$DATE.sql
gzip $BACKUP_DIR/db_$DATE.sql
find $BACKUP_DIR -type f -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/db_$DATE.sql.gz"
BACKUP_EOF

    chmod +x /opt/chengetai/backup.sh
    print_success "Backup script installed"

    # Add to crontab (daily at 2 AM)
    (crontab -l 2>/dev/null | grep -v "backup.sh"; echo "0 2 * * * /opt/chengetai/backup.sh") | crontab -
    print_success "Daily backups scheduled"
}

# Setup monitoring
setup_monitoring() {
    print_header "Setting Up Monitoring"

    cat > /opt/chengetai/monitor.sh << 'MONITOR_EOF'
#!/bin/bash
STATUS=$(curl -s http://localhost:8080/health | grep -c "ok")
if [ $STATUS -eq 0 ]; then
  echo "API DOWN at $(date)" >> /opt/chengetai/logs/alerts.log
  docker-compose -f /opt/chengetai/docker-compose.production.yml restart api
fi
MONITOR_EOF

    chmod +x /opt/chengetai/monitor.sh
    print_success "Monitoring script installed"

    # Add to crontab (every 5 minutes)
    (crontab -l 2>/dev/null | grep -v "monitor.sh"; echo "*/5 * * * * /opt/chengetai/monitor.sh") | crontab -
    print_success "Health monitoring scheduled"
}

# Summary
show_summary() {
    print_header "Deployment Summary"

    echo -e "${GREEN}✓ ChengetAi Learn MVP Deployment Complete!${NC}\n"

    echo "Services running:"
    docker-compose -f docker-compose.production.yml ps

    echo -e "\n${GREEN}Quick Commands:${NC}"
    echo "  View logs:           docker-compose -f docker-compose.production.yml logs -f"
    echo "  Stop services:       docker-compose -f docker-compose.production.yml down"
    echo "  Restart services:    docker-compose -f docker-compose.production.yml restart"
    echo "  Backup database:     /opt/chengetai/backup.sh"
    echo "  View backups:        ls -lah /opt/chengetai/backups/"

    echo -e "\n${GREEN}API Endpoints:${NC}"
    echo "  Health:  curl http://localhost:8080/health"
    echo "  Ready:   curl http://localhost:8080/ready"
    echo "  Version: curl http://localhost:8080/version"

    echo -e "\n${GREEN}Next Steps:${NC}"
    echo "  1. Configure Nginx (see DEPLOY_MVP.md)"
    echo "  2. Set up SSL certificate (Let's Encrypt)"
    echo "  3. Test API endpoints"
    echo "  4. Set up monitoring/alerts"
    echo "  5. Deploy frontend application"

    echo -e "\n${YELLOW}Important:${NC}"
    echo "  - Keep .env file secure and backed up"
    echo "  - Regular database backups are scheduled"
    echo "  - Health monitoring runs every 5 minutes"
    echo "  - Logs are in docker-compose output"

    echo -e "\n${GREEN}Configuration:${NC}"
    echo "  Database User: chengetai"
    echo "  Database Name: chengetai"
    echo "  Redis Port: 6379"
    echo "  API Port: 8080"
}

# Main execution
main() {
    echo -e "\n${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  ChengetAi Learn MVP Deployment Script  ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}\n"

    check_prerequisites
    setup_directories
    generate_secrets
    configure_environment
    build_images
    start_services
    verify_installation
    setup_backups
    setup_monitoring
    show_summary

    print_success "Deployment completed successfully!"
}

# Run main function
main
