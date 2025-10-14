#!/bin/bash
# V-UI - Final Simple Installation
# Guaranteed to work

set -e

echo "========================================
   V-UI Final Installation
========================================"

# Install Docker using official script
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

# Test Docker
/usr/bin/docker ps
/usr/local/bin/docker-compose version

# Clone repository
echo "Cloning V-UI..."
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# Setup environment
echo "Setting up environment..."
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env

# Build with no cache to ensure everything is fresh
echo "Building images (this takes time)..."
/usr/local/bin/docker-compose build --no-cache

# Start services
echo "Starting services..."
/usr/local/bin/docker-compose up -d

# Wait for backend to be ready
echo "Waiting for backend..."
sleep 15
for i in {1..30}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo "✓ Backend is ready!"
        break
    fi
    echo -n "."
    sleep 2
done

# Final message
SERVER_IP=$(hostname -I | awk '{print $1}')
echo ""
echo "========================================
   Installation Complete!
========================================
🌐 Web Panel: http://$SERVER_IP
📚 API Docs: http://$SERVER_IP/docs
👤 Username: admin
🔑 Password: admin

⚠️  IMPORTANT: Change password immediately!
========================================"
