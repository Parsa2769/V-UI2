#!/bin/bash

set -e

echo "3X-UI Modern Installation Script"
echo "================================="

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo "This script should NOT be run as root" 
   exit 1
fi

# Check OS
if [[ ! -f /etc/debian_version ]]; then
    echo "This script is designed for Debian/Ubuntu systems"
    exit 1
fi

echo "Installing dependencies..."

# Update package list
sudo apt-get update

# Install required packages
sudo apt-get install -y \
    curl \
    wget \
    unzip \
    git \
    postgresql \
    postgresql-contrib

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
fi

# Install Docker Compose if not present
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create data directories
echo "Creating data directories..."
mkdir -p data backups config

# Copy environment file
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    
    # Generate random secrets
    JWT_SECRET=$(openssl rand -base64 32)
    JWT_REFRESH_SECRET=$(openssl rand -base64 32)
    DB_PASSWORD=$(openssl rand -base64 24)
    
    sed -i "s/changeme_32_char_random_secret_key_here_minimum/$JWT_SECRET/" .env
    sed -i "s/changeme_32_char_refresh_secret_key_here/$JWT_REFRESH_SECRET/" .env
    sed -i "s/changeme_secure_password/$DB_PASSWORD/" .env
    
    echo "✓ Generated random secrets in .env file"
    echo "⚠️  IMPORTANT: Review and update .env file with your settings"
fi

echo ""
echo "Installation complete!"
echo ""
echo "Next steps:"
echo "1. Review and edit .env file: nano .env"
echo "2. Start services: docker-compose up -d"
echo "3. Check logs: docker-compose logs -f"
echo "4. Access panel: http://localhost:8080"
echo ""
echo "Default credentials:"
echo "  Username: admin"
echo "  Password: admin"
echo "  ⚠️  CHANGE IMMEDIATELY AFTER FIRST LOGIN!"
