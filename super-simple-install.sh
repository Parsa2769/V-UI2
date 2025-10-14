#!/bin/bash
# V-UI - Super Simple Installation
# Uses pre-built simple files

echo "========================================
   V-UI Super Simple Installation
========================================"

# Install Docker (official method)
echo "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh

# Start Docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
echo "Installing Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Test
echo "Testing Docker..."
/usr/bin/docker ps
/usr/local/bin/docker-compose version

# Clone repo
echo "Cloning repository..."
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# Use simple files
echo "Using simple Docker files..."
cp Dockerfile.simple Dockerfile
cp docker-compose.simple.yml docker-compose.yml

# Setup .env
echo "Setting up environment..."
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env

# Build and start
echo "Building images (5-10 minutes)..."
/usr/local/bin/docker-compose build

echo "Starting services..."
/usr/local/bin/docker-compose up -d

# Wait for backend
sleep 15
echo "Waiting for backend..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo "✓ Backend is ready!"
        break
    fi
    echo -n "."
    sleep 2
done

echo ""
echo "========================================
   Installation Complete!
========================================
🌐 Web Panel: http://$(hostname -I | awk '{print $1}')"
echo "👤 Username: admin"
echo "🔑 Password: admin"
echo ""
echo "⚠️  CHANGE PASSWORD IMMEDIATELY!"
echo "========================================"
