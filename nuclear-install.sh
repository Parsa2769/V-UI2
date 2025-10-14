#!/bin/bash
# V-UI - Nuclear Installation (Everything from scratch)
# Destroys everything and rebuilds from zero

set -e

echo "========================================
   V-UI Nuclear Installation
   Everything from scratch
========================================"

# Stop and remove everything
echo "Cleaning up previous installation..."
docker compose down -v 2>/dev/null || true
docker system prune -af 2>/dev/null || true
rm -rf V-UI2 2>/dev/null || true

# Install Docker from scratch
echo "Installing Docker from scratch..."
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

# Clone fresh repository
echo "Cloning fresh repository..."
git clone https://github.com/Parsa2769/V-UI2.git
cd V-UI2

# Fix any TypeScript issues
echo "Fixing TypeScript issues..."
sed -i 's/import { render, screen }/import { render }/' frontend/src/App.test.tsx
sed -i 's/getActionIcon(log.action)/getActionIcon()/' frontend/src/pages/AuditLogsPage.tsx

# Setup environment
echo "Setting up environment..."
cp .env.example .env
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
DB_PASSWORD=$(openssl rand -base64 24 | tr -d '\n')
sed -i "s|changeme_32_char_random_secret_key_here_minimum|$JWT_SECRET|g" .env
sed -i "s|changeme_secure_password|$DB_PASSWORD|g" .env

# Build with nuclear option (no cache, fresh everything)
echo "Building from scratch (this takes 10-15 minutes)..."
/usr/local/bin/docker-compose build --no-cache --pull

# Start services
echo "Starting services..."
/usr/local/bin/docker-compose up -d

# Wait for backend
sleep 20
echo "Waiting for backend to be fully ready..."
for i in {1..45}; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo "✓ Backend is ready!"
        break
    fi
    echo -n "."
    sleep 2
done

# Final check
echo ""
echo "========================================
   Nuclear Installation Complete!
========================================
🌐 Web Panel: http://$(hostname -I | awk '{print $1}')"
echo "📚 API Docs: http://$(hostname -I | awk '{print $1}')/docs"
echo "👤 Username: admin"
echo "🔑 Password: admin"
echo ""
echo "⚠️  CHANGE PASSWORD IMMEDIATELY!"
echo "========================================"
echo ""
echo "Useful commands:"
echo "  Check status: docker compose ps"
echo "  View logs: docker compose logs -f"
echo "  Restart: docker compose restart"
echo "  Stop: docker compose down"
