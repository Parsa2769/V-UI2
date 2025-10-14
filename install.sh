#!/bin/bash

# V-UI - Complete Installation Script
# Auto-installs Docker if not found

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   V-UI Installation Script${NC}"
echo -e "${GREEN}   Modern Xray Management Panel${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${YELLOW}⚠️  Not running as root. Using sudo.${NC}"
   SUDO='sudo'
else
   echo -e "${GREEN}✓ Running as root${NC}"
   SUDO=''
fi
echo ""

# Detect OS
if [[ ! -f /etc/debian_version ]] && [[ ! -f /etc/lsb-release ]]; then
    echo -e "${RED}❌ This script is designed for Debian/Ubuntu systems${NC}"
    exit 1
fi

OS_NAME=$(lsb_release -si)
OS_VERSION=$(lsb_release -sr)
echo -e "Detected OS: ${GREEN}$OS_NAME $OS_VERSION${NC}"
echo ""

# Install Docker if not present
echo -e "${BLUE}[1/7]${NC} ${YELLOW}Checking Docker installation...${NC}"

DOCKER_INSTALLED=false
DOCKER_BIN=""

if [ -x "/usr/bin/docker" ]; then
    DOCKER_BIN="/usr/bin/docker"
    DOCKER_INSTALLED=true
    echo -e "   Found Docker at: ${DOCKER_BIN}"
elif [ -x "/usr/local/bin/docker" ]; then
    DOCKER_BIN="/usr/local/bin/docker"
    DOCKER_INSTALLED=true
    echo -e "   Found Docker at: ${DOCKER_BIN}"
elif command -v docker &> /dev/null; then
    DOCKER_BIN="docker"
    DOCKER_INSTALLED=true
    echo -e "   Docker found in PATH"
else
    echo -e "   ${YELLOW}Docker not found. Installing...${NC}"

    # Use official Docker installation script
    echo -e "   Downloading official Docker installer..."
    curl -fsSL https://get.docker.com -o get-docker.sh

    echo -e "   Installing Docker (this may take 2-3 minutes)..."
    sh get-docker.sh

    # Clean up
    rm get-docker.sh

    # Find Docker after installation
    if [ -x "/usr/bin/docker" ]; then
        DOCKER_BIN="/usr/bin/docker"
        DOCKER_INSTALLED=true
        echo -e "${GREEN}✓ Docker installed successfully${NC}"
    else
        echo -e "${RED}❌ Failed to install Docker${NC}"
        echo "   Please install manually: https://docs.docker.com/engine/install/ubuntu/"
        exit 1
    fi
fi

# Start Docker service
echo -e "   Starting Docker service..."
$SUDO systemctl daemon-reload 2>/dev/null || true
$SUDO systemctl start docker 2>/dev/null || true
$SUDO systemctl enable docker 2>/dev/null || true

# Wait for Docker daemon
for i in {1..10}; do
    if $DOCKER_BIN ps &> /dev/null 2>&1; then
        echo -e "${GREEN}✓ Docker is working${NC}"
        break
    fi
    sleep 1
done

# Verify Docker works
if ! $DOCKER_BIN ps &> /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Docker daemon not responding${NC}"
    echo -e "   Trying direct daemon start..."
    $SUDO dockerd &> /dev/null &
    sleep 3

    if $DOCKER_BIN ps &> /dev/null 2>&1; then
        echo -e "${GREEN}✓ Docker daemon started${NC}"
    else
        echo -e "${RED}❌ Docker daemon failed to start${NC}"
        echo "   Manual fix: sudo systemctl start docker"
        exit 1
    fi
fi

# Install Docker Compose if needed
echo ""
echo -e "${BLUE}[2/7]${NC} ${YELLOW}Checking Docker Compose...${NC}"

COMPOSE_CMD=""
if $DOCKER_BIN compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="$DOCKER_BIN compose"
    echo -e "   Found Docker Compose plugin"
elif [ -x "/usr/local/bin/docker-compose" ]; then
    COMPOSE_CMD="/usr/local/bin/docker-compose"
    echo -e "   Found Docker Compose standalone"
elif [ -x "/usr/bin/docker-compose" ]; then
    COMPOSE_CMD="/usr/bin/docker-compose"
    echo -e "   Found Docker Compose"
else
    echo -e "   ${YELLOW}Installing Docker Compose...${NC}"

    # Install Docker Compose
    DOCKER_COMPOSE_VERSION="2.24.5"
    $SUDO curl -L "https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    $SUDO chmod +x /usr/local/bin/docker-compose

    COMPOSE_CMD="/usr/local/bin/docker-compose"
    echo -e "${GREEN}✓ Docker Compose installed${NC}"
fi

echo -e "   Using compose command: ${COMPOSE_CMD}"

# Clone repository
echo ""
echo -e "${BLUE}[3/7]${NC} ${YELLOW}Cloning V-UI repository...${NC}"

if [ -d "V-UI2" ]; then
    echo -e "   Directory exists, pulling latest..."
    cd V-UI2
    git pull
else
    git clone https://github.com/Parsa2769/V-UI2.git
    cd V-UI2
    echo -e "${GREEN}✓ Repository cloned${NC}"
fi

# Setup environment
echo ""
echo -e "${BLUE}[4/7]${NC} ${YELLOW}Setting up environment...${NC}"

if [ ! -f .env ]; then
    echo -e "   Creating .env file..."
    cp .env.example .env

    echo -e "   Generating secure secrets..."
    JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
    JWT_REFRESH_SECRET=$(openssl rand -base64 32 | tr -d '\n')
    DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')

    sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
    sed -i "s|changeme_32_char_refresh_secret_key_here|$JWT_REFRESH_SECRET|g" .env
    sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env

    echo -e "${GREEN}✓ Environment configured${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Build images
echo ""
echo -e "${BLUE}[5/7]${NC} ${YELLOW}Building Docker images...${NC}"
echo "   This takes 5-10 minutes (downloading Go, Node.js, building)..."

if $COMPOSE_CMD build --no-cache; then
    echo -e "${GREEN}✓ Images built successfully${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    echo "   Check logs above for errors"
    exit 1
fi

# Start services
echo ""
echo -e "${BLUE}[6/7]${NC} ${YELLOW}Starting services...${NC}"

if $COMPOSE_CMD up -d; then
    echo -e "${GREEN}✓ Services started${NC}"
else
    echo -e "${RED}❌ Failed to start services${NC}"
    echo "   Run: $COMPOSE_CMD logs"
    exit 1
fi

# Wait for backend
echo ""
echo -e "${BLUE}[7/7]${NC} ${YELLOW}Waiting for backend...${NC}"
sleep 10

for i in {1..30}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Backend is healthy${NC}"
        break
    fi
    echo -n "."
    sleep 2
done
echo ""

# Get server IP
SERVER_IP=$(hostname -I | awk '{print $1}' || echo "localhost")

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ Installation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "🌐 ${BLUE}Web Panel:${NC} http://$SERVER_IP"
echo -e "📚 ${BLUE}API Docs:${NC}  http://$SERVER_IP/docs"
echo ""
echo -e "👤 ${YELLOW}Default Login:${NC}"
echo -e "   Username: ${GREEN}admin${NC}"
echo -e "   Password: ${GREEN}admin${NC}"
echo ""
echo -e "${RED}⚠️  CHANGE PASSWORD IMMEDIATELY!${NC}"
echo ""
echo -e "${BLUE}Useful Commands:${NC}"
echo -e "  View logs:    ${YELLOW}$COMPOSE_CMD logs -f${NC}"
echo -e "  Stop:         ${YELLOW}$COMPOSE_CMD down${NC}"
echo -e "  Restart:      ${YELLOW}$COMPOSE_CMD restart${NC}"
echo ""
