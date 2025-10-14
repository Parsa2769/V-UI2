#!/bin/bash
# V-UI - Simple Installation Script (Fixed)
# Installs Docker, clones repo, and runs everything

set -e

echo "========================================
   V-UI Installation Script
   Fixed Version
========================================"

# Install Docker using official script
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

# Clone or update repo
echo "Cloning V-UI repository..."
if [ -d "V-UI2" ]; then
    cd V-UI2
    git pull
else
    git clone https://github.com/Parsa2769/V-UI2.git
    cd V-UI2
fi

# Setup .env
echo "Setting up environment..."
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env

# Build and start
echo "Building Docker images (5-10 minutes)..."
/usr/local/bin/docker-compose build

echo "Starting services..."
/usr/local/bin/docker-compose up -d

# Wait for backend
sleep 15
echo "Waiting for backend to be ready..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo "✓ Backend is ready!"
        break
    fi
    echo -n "."
    sleep 2
done

# Final status
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
