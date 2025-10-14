#!/bin/bash

<<<<<<< HEAD
# V-UI - Modern Xray Management Panel
# Installation Script for Ubuntu/Debian
# Tested on Ubuntu 20.04, 22.04, 24.04
=======
# V-UI - Complete Installation Script
# Auto-installs Docker if not found
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
<<<<<<< HEAD
=======
BLUE='\033[0;34m'
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   V-UI Installation Script${NC}"
echo -e "${GREEN}   Modern Xray Management Panel${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
<<<<<<< HEAD
   echo -e "${YELLOW}⚠️  Not running as root. Some commands may require sudo.${NC}"
=======
   echo -e "${YELLOW}⚠️  Not running as root. Using sudo.${NC}"
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281
   SUDO='sudo'
else
   echo -e "${GREEN}✓ Running as root${NC}"
   SUDO=''
fi
<<<<<<< HEAD

# Check OS
=======
echo ""

# Detect OS
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281
if [[ ! -f /etc/debian_version ]] && [[ ! -f /etc/lsb-release ]]; then
    echo -e "${RED}❌ This script is designed for Debian/Ubuntu systems${NC}"
    exit 1
fi

<<<<<<< HEAD
# Get OS info
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
    VER=$VERSION_ID
    echo -e "Detected OS: ${GREEN}$OS $VER${NC}"
fi

echo ""
echo -e "${YELLOW}📦 Installing system dependencies...${NC}"

# Update package list
echo -e "   Updating package list..."
$SUDO apt-get update -qq 2>&1 | grep -v "^Get:" | grep -v "^Hit:" || true

# Install required packages
echo -e "   Installing curl, wget, git, openssl..."
$SUDO DEBIAN_FRONTEND=noninteractive apt-get install -y \
    curl \
    wget \
    unzip \
    git \
    ca-certificates \
    gnupg \
    lsb-release \
    openssl 2>&1 | grep -E "(Setting up|Unpacking|already)" || true

echo -e "${GREEN}✓ System dependencies installed${NC}"

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    echo ""
    echo -e "${YELLOW}🐳 Installing Docker...${NC}"
    
    # Add Docker's official GPG key
    echo -e "   Adding Docker GPG key..."
    $SUDO mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | $SUDO gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    
    # Set up the repository
    echo -e "   Setting up Docker repository..."
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | $SUDO tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Install Docker Engine
    echo -e "   Installing Docker Engine (this may take a few minutes)..."
    $SUDO apt-get update -qq 2>&1 | grep -v "^Get:" | grep -v "^Hit:" || true
    $SUDO DEBIAN_FRONTEND=noninteractive apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin 2>&1 | grep -E "(Setting up|Unpacking|already)" || true
    
    # Add current user to docker group (if not root)
    if [[ $EUID -ne 0 ]]; then
        $SUDO usermod -aG docker $USER
    fi
    
    echo -e "${GREEN}✓ Docker installed successfully${NC}"
    echo -e "${YELLOW}⚠️  Please log out and log back in for docker group changes to take effect${NC}"
else
    echo -e "${GREEN}✓ Docker is already installed${NC}"
fi

# Check Docker Compose
if ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose plugin not found${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker Compose is available${NC}"

# Create data directories
echo ""
echo -e "${YELLOW}📁 Creating data directories...${NC}"
mkdir -p data backups config nginx/ssl

echo -e "${GREEN}✓ Directories created${NC}"

# Copy environment file
if [ ! -f .env ]; then
    echo ""
    echo -e "${YELLOW}🔧 Creating .env file...${NC}"
    
    if [ ! -f .env.example ]; then
        echo -e "${RED}❌ .env.example not found${NC}"
        exit 1
    fi
    
    cp .env.example .env
    
    # Generate random secrets
    echo -e "${YELLOW}🔐 Generating secure random secrets...${NC}"
    JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
    JWT_REFRESH_SECRET=$(openssl rand -base64 32 | tr -d '\n')
    DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')
    
    # Update .env file
    sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
    sed -i "s|changeme_32_char_refresh_secret_key_here|$JWT_REFRESH_SECRET|g" .env
    sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env
    
    echo -e "${GREEN}✓ Generated random secrets in .env file${NC}"
    echo -e "${YELLOW}⚠️  IMPORTANT: Review and update .env file with your settings${NC}"
=======
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
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

<<<<<<< HEAD
# Build and start services
echo ""
echo -e "${YELLOW}🏗️  Building Docker images...${NC}"
echo "   This may take 5-10 minutes on first run..."
echo "   (Downloading Go, Node.js, and building both backend and frontend)"

if docker compose build --no-cache 2>&1 | grep -E "(Step|Successfully)" || docker compose build --no-cache; then
    echo -e "${GREEN}✓ Docker images built successfully${NC}"
else
    echo -e "${RED}❌ Failed to build Docker images${NC}"
    echo "   Check logs above for details"
    exit 1
fi

echo ""
echo -e "${YELLOW}🚀 Starting services...${NC}"

if docker compose up -d 2>&1; then
    echo -e "${GREEN}✓ Services started successfully${NC}"
else
    echo -e "${RED}❌ Failed to start services${NC}"
    echo "   Run: docker compose logs"
    exit 1
fi

# Wait for services to be ready
echo ""
echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
sleep 10

# Check if backend is healthy
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
=======
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
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Backend is healthy${NC}"
        break
    fi
<<<<<<< HEAD
    
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo -n "."
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo ""
    echo -e "${RED}❌ Backend health check failed${NC}"
    echo "   Check logs: docker compose logs backend"
    exit 1
fi

# Get server IP
SERVER_IP=$(hostname -I | awk '{print $1}')
=======
    echo -n "."
    sleep 2
done
echo ""

# Get server IP
SERVER_IP=$(hostname -I | awk '{print $1}' || echo "localhost")
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ Installation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
<<<<<<< HEAD
echo -e "Access URLs:"
echo -e "  Web Panel: ${GREEN}http://$SERVER_IP${NC} or ${GREEN}http://localhost${NC}"
echo -e "  API Docs:  ${GREEN}http://$SERVER_IP/docs${NC}"
echo -e "  Metrics:   ${GREEN}http://$SERVER_IP/metrics${NC}"
echo ""
echo -e "Default Credentials:"
echo -e "  Username: ${YELLOW}admin${NC}"
echo -e "  Password: ${YELLOW}admin${NC}"
echo ""
echo -e "${RED}⚠️  CRITICAL: Change default password immediately after first login!${NC}"
echo ""
echo "Useful Commands:"
echo "  View logs:    docker compose logs -f"
echo "  Stop:         docker compose down"
echo "  Restart:      docker compose restart"
echo "  Update:       git pull && docker compose up -d --build"
echo "  Backup:       ./scripts/backup.sh"
echo ""
echo -e "${YELLOW}📚 Documentation: https://github.com/Parsa2769/V-UI2${NC}"
=======
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
>>>>>>> 1a83f2e275b89d9fb23847d0787f205c8d23a281
echo ""
