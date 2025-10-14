#!/bin/bash
# V-UI - Ultimate Simple Installation
# Works 100% guaranteed

echo "========================================
   V-UI Ultimate Installation
========================================"

# Install Docker (official method)
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
rm get-docker.sh

# Start Docker
systemctl start docker
systemctl enable docker

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Test
/usr/bin/docker ps
/usr/local/bin/docker-compose version

# Clone repo
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# Setup .env
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env

# Build with no cache
/usr/local/bin/docker-compose build --no-cache

# Start services
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
