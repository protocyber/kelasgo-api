#!/bin/bash

# Database Migration Helper Script
# This script handles all database migration operations

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if yq is installed
check_yq() {
    if ! command -v yq >/dev/null 2>&1; then
        echo -e "${RED}❌ yq not found. Install with:${NC}"
        echo -e "${YELLOW}  macOS: brew install yq${NC}"
        echo -e "${YELLOW}  Linux: Check https://github.com/mikefarah/yq${NC}"
        exit 1
    fi
}

# Check if config.yaml exists
check_config() {
    if [ ! -f "configs/config.yaml" ]; then
        echo -e "${RED}❌ configs/config.yaml not found${NC}"
        echo -e "${YELLOW}Please create it from configs/config.example.yaml:${NC}"
        echo -e "  cp configs/config.example.yaml configs/config.yaml"
        echo -e "  # Edit and update database credentials"
        exit 1
    fi
}

# Extract database configuration
get_db_config() {
    check_yq
    check_config
    
    DB_HOST=$(yq '.db.pg.write.host' configs/config.yaml 2>/dev/null || echo "localhost")
    DB_PORT=$(yq '.db.pg.write.port' configs/config.yaml 2>/dev/null || echo "5432")
    DB_NAME=$(yq '.db.pg.write.name' configs/config.yaml 2>/dev/null || echo "kelasgo")
    DB_USER=$(yq '.db.pg.write.user' configs/config.yaml 2>/dev/null || echo "postgres")
    DB_PASSWORD=$(yq '.db.pg.write.password' configs/config.yaml 2>/dev/null || echo "")
    DB_SSLMODE=$(yq '.db.pg.write.sslmode' configs/config.yaml 2>/dev/null || echo "disable")
    
    DB_CONN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
}

# Show database configuration
show_config() {
    get_db_config
    echo -e "${BLUE}📊 Database Configuration:${NC}"
    echo -e "   Host: ${DB_HOST}"
    echo -e "   Port: ${DB_PORT}"
    echo -e "   Database: ${DB_NAME}"
    echo -e "   User: ${DB_USER}"
    echo -e "   SSL Mode: ${DB_SSLMODE}"
}

# Create new migration
create_migration() {
    local migration_name="$1"
    if [ -z "$migration_name" ]; then
        echo -e "${RED}❌ Migration name required${NC}"
        echo -e "${YELLOW}Usage: $0 create <migration_name>${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}📝 Creating migration: ${migration_name}${NC}"
    migrate create -ext sql -dir ./migrations/postgres "$migration_name"
}

# Run migrations up
migrate_up() {
    local steps="${1:-1}"
    get_db_config
    echo -e "${BLUE}⬆️  Running migrations up (${steps} steps)${NC}"
    migrate -path ./migrations/postgres -database "$DB_CONN" up "$steps"
}

# Run migrations down
migrate_down() {
    local steps="${1:-1}"
    get_db_config
    echo -e "${BLUE}⬇️  Running migrations down (${steps} steps)${NC}"
    migrate -path ./migrations/postgres -database "$DB_CONN" down "$steps"
}

# Force migration version
migrate_force() {
    local version="$1"
    if [ -z "$version" ]; then
        echo -e "${RED}❌ Migration version required${NC}"
        echo -e "${YELLOW}Usage: $0 force <version>${NC}"
        exit 1
    fi
    
    get_db_config
    echo -e "${BLUE}🔧 Forcing migration version: ${version}${NC}"
    migrate -path ./migrations/postgres -database "$DB_CONN" force "$version"
}

# Show migration version
migrate_version() {
    get_db_config
    echo -e "${BLUE}📋 Current migration version:${NC}"
    migrate -path ./migrations/postgres -database "$DB_CONN" version
}

# Drop database (dangerous!)
migrate_drop() {
    get_db_config
    echo -e "${RED}⚠️  WARNING: This will drop all data!${NC}"
    read -p "Are you sure? Type 'yes' to confirm: " confirm
    if [ "$confirm" = "yes" ]; then
        migrate -path ./migrations/postgres -database "$DB_CONN" drop
    else
        echo -e "${YELLOW}Operation cancelled${NC}"
    fi
}

# Main command dispatcher
case "$1" in
    "config")
        show_config
        ;;
    "create")
        create_migration "$2"
        ;;
    "up")
        migrate_up "$2"
        ;;
    "down")
        migrate_down "$2"
        ;;
    "force")
        migrate_force "$2"
        ;;
    "version")
        migrate_version
        ;;
    "drop")
        migrate_drop
        ;;
    *)
        echo -e "${BLUE}📖 Database Migration Helper${NC}"
        echo -e "${YELLOW}Usage: $0 <command> [args]${NC}"
        echo ""
        echo -e "${BLUE}Commands:${NC}"
        echo -e "  config              - Show database configuration"
        echo -e "  create <name>       - Create new migration"
        echo -e "  up [steps]         - Run migrations up (default: 1)"
        echo -e "  down [steps]       - Run migrations down (default: 1)"
        echo -e "  force <version>    - Force migration version"
        echo -e "  version            - Show current migration version"
        echo -e "  drop               - Drop database (WARNING: destructive)"
        exit 1
        ;;
esac
