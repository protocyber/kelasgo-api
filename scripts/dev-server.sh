#!/bin/bash

# OS-Aware Development Server Script
# Detects OS and starts the appropriate development environment

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect operating system
detect_os() {
    case "$(uname -s)" in
        "Darwin")
            OS="macos"
            OS_NAME="macOS"
            OS_ICON="🖥️"
            ;;
        "Linux")
            OS="linux"
            OS_NAME="Linux"
            OS_ICON="🐧"
            ;;
        "MINGW"* | "MSYS"* | "CYGWIN"*)
            OS="windows"
            OS_NAME="Windows"
            OS_ICON="🪟"
            ;;
        *)
            OS="unknown"
            OS_NAME="Unknown"
            OS_ICON="❓"
            ;;
    esac
}

# Check if Air is installed
check_air() {
    local air_path="$(go env GOPATH)/bin/air"
    if [ ! -f "$air_path" ]; then
        echo -e "${RED}❌ Air not found at: $air_path${NC}"
        echo -e "${YELLOW}Install Air with: go install github.com/cosmtrek/air@latest${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Air found: $air_path${NC}"
}

# Setup macOS specific environment
setup_macos() {
    echo -e "${BLUE}🔧 Setting up macOS environment...${NC}"
    
    # Check if symlinks exist
    if [ ! -L /usr/local/include/leptonica ] || [ ! -L /usr/local/include/tesseract ]; then
        echo -e "${YELLOW}🔗 Running macOS symlink setup...${NC}"
        bash scripts/macos-setup.sh
    else
        echo -e "${GREEN}✅ macOS symlinks already set up${NC}"
    fi
    
    # Set macOS specific environment variables
    export CGO_CFLAGS="-I/opt/homebrew/opt/leptonica/include -I/opt/homebrew/opt/tesseract/include"
    export CGO_LDFLAGS="-L/opt/homebrew/opt/leptonica/lib -L/opt/homebrew/opt/tesseract/lib"
    export PKG_CONFIG_PATH="/opt/homebrew/opt/leptonica/lib/pkgconfig:/opt/homebrew/opt/tesseract/lib/pkgconfig"
    export DYLD_LIBRARY_PATH="/opt/homebrew/opt/leptonica/lib:/opt/homebrew/opt/tesseract/lib"
}

# Start development server
start_dev_server() {
    echo -e "${BLUE}🌐 Detecting operating system...${NC}"
    detect_os
    echo -e "${GREEN}${OS_ICON} ${OS_NAME} detected${NC}"
    
    # OS-specific setup
    case "$OS" in
        "macos")
            setup_macos
            ;;
        "linux"|"windows")
            echo -e "${BLUE}🐧 Using standard Linux/Windows environment${NC}"
            ;;
        *)
            echo -e "${YELLOW}⚠️  Unknown OS, using default environment${NC}"
            ;;
    esac
    
    # Check dependencies
    check_air
    
    # Start Air
    echo -e "${GREEN}🚀 Starting development server with Air...${NC}"
    "$(go env GOPATH)/bin/air"
}

# Check configuration
check_config() {
    if [ -f "configs/config.yaml" ]; then
        echo -e "${GREEN}✅ configs/config.yaml found${NC}"
    else
        echo -e "${RED}❌ configs/config.yaml not found${NC}"
        echo -e "${YELLOW}Please create it from configs/config.example.yaml:${NC}"
        echo -e "  cp configs/config.example.yaml configs/config.yaml"
        echo -e "  # Edit and update your settings"
        exit 1
    fi
}

# Show environment info
show_env() {
    detect_os
    echo -e "${BLUE}📋 Development Environment:${NC}"
    echo -e "  OS: ${OS_ICON} ${OS_NAME}"
    echo -e "  Go Version: $(go version | cut -d' ' -f3)"
    echo -e "  GOPATH: $(go env GOPATH)"
    echo -e "  Air: $([ -f "$(go env GOPATH)/bin/air" ] && echo "✅ Installed" || echo "❌ Not installed")"
    if [ -f "configs/config.yaml" ]; then
        echo -e "  Config: ✅ configs/config.yaml found"
    else
        echo -e "  Config: ❌ configs/config.yaml missing"
    fi
}

# Main command dispatcher
case "${1:-start}" in
    "start")
        check_config
        start_dev_server
        ;;
    "env")
        show_env
        ;;
    "check")
        check_config
        check_air
        echo -e "${GREEN}✅ All checks passed${NC}"
        ;;
    *)
        echo -e "${BLUE}📖 Development Server Helper${NC}"
        echo -e "${YELLOW}Usage: $0 [command]${NC}"
        echo ""
        echo -e "${BLUE}Commands:${NC}"
        echo -e "  start (default)    - Start development server"
        echo -e "  env               - Show environment information"
        echo -e "  check             - Check dependencies and configuration"
        exit 1
        ;;
esac
