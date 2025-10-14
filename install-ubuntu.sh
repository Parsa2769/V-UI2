#!/bin/bash

# V-UI - Modern Xray Management Panel
# Installation Script for Ubuntu/Debian
# Tested on Ubuntu 20.04, 22.04, 24.04

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   V-UI Installation Script${NC}"
echo -e "${GREEN}   Modern Xray Management Panel${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${YELLOW}⚠️  Not running as root. Some commands may require sudo.${NC}"
   SUDO='sudo'
else
   echo -e "${GREEN}✓ Running as root${NC}"
   SUDO=''
fi

# Check OS
if [[ ! -f /etc/debian_version ]] && [[ ! -f /etc/lsb-release ]]; then
    echo -e "${RED}❌ This script is designed for Debian/Ubuntu systems${NC}"
    exit 1
fi

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
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

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
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Backend is healthy${NC}"
        break
    fi
    
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

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ Installation Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
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
echo ""
