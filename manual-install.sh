#!/bin/bash
# V-UI - Manual Installation Script
# Step by step, no automation issues

echo "========================================
   V-UI Manual Installation
   Step by Step
========================================"

echo "Step 1: Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh

echo "Step 2: Starting Docker..."
systemctl start docker
systemctl enable docker

echo "Step 3: Installing Docker Compose..."
curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo "Step 4: Testing Docker..."
/usr/bin/docker ps
/usr/local/bin/docker-compose version

echo "Step 5: Cloning repository..."
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

echo "Step 6: Setting up environment..."
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env

echo "Step 7: Building images..."
echo "This will take 5-10 minutes..."
/usr/local/bin/docker-compose build

echo "Step 8: Starting services..."
/usr/local/bin/docker-compose up -d

echo "Step 9: Waiting for backend..."
sleep 15
echo "Checking backend health..."
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
