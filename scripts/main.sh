#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

detect_os() {
    case "$(uname -s)" in
        Linux*)
            if [[ -f /etc/os-release ]]; then
                . /etc/os-release
                echo "Linux-$ID"
            else
                echo "Linux"
            fi
            ;;
        Darwin*)    echo "macOS";;
        CYGWIN*)    echo "Windows";;
        MINGW*)     echo "Windows";;
        MSYS*)      echo "Windows";;
        *)          echo "Unknown";;
    esac
}

# Check if tool exists
check_tool() {
    local toolName=$1
    if command -v "$toolName" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ ${toolName} is installed${NC}"
        return 0
    else
        echo -e "${RED}❌ ${toolName} is not installed${NC}"
        return 1
    fi
}

# Install tool based on OS
install_tool() {
    local tool=$1
    local os=$2

    echo -e "${YELLOW}🔧 Installing $tool...${NC}"

    case $os in
        macOS)
            if command -v brew >/dev/null 2>&1; then
                brew install "$tool"
            else
                echo -e "${RED}❌ Homebrew not found. Please install it first.${NC}"
                return 1
            fi
            ;;
        Linux-ubuntu|Linux-debian)
            sudo apt update && sudo apt install -y "$tool"
            ;;
        Linux-fedora|Linux-rhel|Linux-centos)
            sudo dnf install -y "$tool" || sudo yum install -y "$tool"
            ;;
        Linux-arch)
            sudo pacman -S "$tool"
            ;;
        Windows)
            echo -e "${YELLOW}⚠️ Please install $tool manually on Windows${NC}"
            return 1
            ;;
        *)
            echo -e "${RED}❌ Unsupported OS: $os${NC}"
            return 1
            ;;
    esac
}

# Main function
main() {
    echo "=== System Information ==="
    OS=$(detect_os)
    echo "Detected OS: $OS"
    echo "Architecture: $(uname -m)"
    echo ""

    echo "=== Tool Check ==="

    REQUIRED_TOOLS=("git"  "go")

    for tool in "${REQUIRED_TOOLS[@]}"; do
        if ! check_tool "$tool"; then
            read -p "Do you want to install $tool? (y/n): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                install_tool "$tool" "$OS"
            fi
        fi
    done

    echo ""
    echo "=== Version Information ==="
    for tool in "${REQUIRED_TOOLS[@]}"; do
        if command -v "$tool" >/dev/null 2>&1; then
            case $tool in
                docker)
                    echo "$tool: $(docker --version)"
                    ;;
                kubectl)
                    echo "$tool: $(kubectl version --client --short 2>/dev/null || echo 'Not connected to cluster')"
                    ;;
                terraform)
                    echo "$tool: $(terraform version -json 2>/dev/null | grep -o '"terraform_version":"[^"]*' | cut -d'"' -f4 || terraform version)"
                    ;;
                go)
                    echo "$tool: $(go version)"
                    ;;
                *)
                    echo "$tool: $(command -v "$tool")"
                    ;;
            esac
        fi
    done
}

main "$@"