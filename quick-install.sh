#!/bin/bash
# V-UI Quick Install - No frills version

set -e

echo "========================================
   V-UI Quick Installation
========================================"
echo ""

# Install Docker using get.docker.com
echo "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh

# Start Docker
systemctl start docker
systemctl enable docker

# Wait for Docker
sleep 5

# Install Docker Compose
echo "Installing Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Test
/usr/bin/docker ps
/usr/local/bin/docker-compose version

# Clone repo
echo "Cloning repository..."
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# Setup .env
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env

# Build
echo "Building (this takes 5-10 minutes)..."
/usr/local/bin/docker-compose build

# Start
echo "Starting services..."
/usr/local/bin/docker-compose up -d

# Wait
sleep 10

# Done
SERVER_IP=$(hostname -I | awk '{print $1}')
echo ""
echo "========================================
   Installation Complete!
========================================
Web Panel: http://$SERVER_IP
Username: admin
Password: admin

⚠️  CHANGE PASSWORD IMMEDIATELY!
========================================"
