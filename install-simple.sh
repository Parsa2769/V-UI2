#!/bin/bash

# V-UI - Simple Installation Script (Verbose)
# For debugging installation issues

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

# 1. Update system
echo -e "${BLUE}[1/7]${NC} ${YELLOW}Updating package list...${NC}"
$SUDO apt-get update -y
echo -e "${GREEN}✓ Package list updated${NC}"
echo ""

# 2. Install dependencies
echo -e "${BLUE}[2/7]${NC} ${YELLOW}Installing dependencies...${NC}"
$SUDO apt-get install -y curl wget git unzip ca-certificates gnupg lsb-release openssl
echo -e "${GREEN}✓ Dependencies installed${NC}"
echo ""

# 3. Install Docker
if ! command -v docker &> /dev/null; then
    echo -e "${BLUE}[3/7]${NC} ${YELLOW}Installing Docker...${NC}"
    
    $SUDO mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | $SUDO gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    
    # Get Ubuntu codename (fallback to noble for unsupported versions)
    UBUNTU_CODENAME=$(lsb_release -cs)
    if [[ "$UBUNTU_CODENAME" == "plucky" ]] || [[ ! "$UBUNTU_CODENAME" =~ ^(focal|jammy|noble)$ ]]; then
        echo -e "   ${YELLOW}Using Ubuntu noble (24.04) repository${NC}"
        UBUNTU_CODENAME="noble"
    fi
    
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $UBUNTU_CODENAME stable" | $SUDO tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    $SUDO apt-get update -y
    $SUDO apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    echo -e "${GREEN}✓ Docker installed${NC}"
else
    echo -e "${BLUE}[3/7]${NC} ${GREEN}✓ Docker already installed${NC}"
fi

# Check Docker Compose
if ! docker compose version &> /dev/null; then
    echo -e "${YELLOW}Installing Docker Compose manually...${NC}"
    DOCKER_COMPOSE_VERSION="2.24.5"
    $SUDO curl -L "https://github.com/docker/compose/releases/download/v${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    $SUDO chmod +x /usr/local/bin/docker-compose
    $SUDO ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    echo -e "${GREEN}✓ Docker Compose installed${NC}"
fi
echo ""

# 4. Clone repository
echo -e "${BLUE}[4/7]${NC} ${YELLOW}Cloning V-UI repository...${NC}"
if [ -d "V-UI2" ]; then
    echo -e "${YELLOW}   Directory exists, pulling latest changes...${NC}"
    cd V-UI2
    git pull
else
    git clone https://github.com/Parsa2769/V-UI2.git
    cd V-UI2
fi
echo -e "${GREEN}✓ Repository ready${NC}"
echo ""

# 5. Setup environment
echo -e "${BLUE}[5/7]${NC} ${YELLOW}Setting up environment...${NC}"
if [ ! -f .env ]; then
    cp .env.example .env
    
    JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
    JWT_REFRESH_SECRET=$(openssl rand -base64 32 | tr -d '\n')
    DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')
    
    sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
    sed -i "s|changeme_32_char_refresh_secret_key_here|$JWT_REFRESH_SECRET|g" .env
    sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env
    
    echo -e "${GREEN}✓ Environment configured with secure secrets${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi
echo ""

# 6. Build
echo -e "${BLUE}[6/7]${NC} ${YELLOW}Building Docker images...${NC}"
echo -e "${YELLOW}   This will take 5-10 minutes...${NC}"
docker compose build
echo -e "${GREEN}✓ Images built${NC}"
echo ""

# 7. Start
echo -e "${BLUE}[7/7]${NC} ${YELLOW}Starting services...${NC}"
docker compose up -d
echo -e "${GREEN}✓ Services started${NC}"
echo ""

# Wait for backend
echo -e "${YELLOW}⏳ Waiting for backend to be ready...${NC}"
sleep 5

for i in {1..30}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Backend is healthy${NC}"
        break
    fi
    echo -n "."
    sleep 2
done
echo ""

# Get IP
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
echo -e "  View logs:    ${YELLOW}docker compose logs -f${NC}"
echo -e "  Stop:         ${YELLOW}docker compose down${NC}"
echo -e "  Restart:      ${YELLOW}docker compose restart${NC}"
echo ""
