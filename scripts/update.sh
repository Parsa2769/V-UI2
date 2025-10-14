#!/bin/bash

set -e

echo "Updating 3X-UI Modern..."

# Backup before update
echo "Creating backup..."
./scripts/backup.sh

# Pull latest changes
echo "Pulling latest changes..."
git pull

# Rebuild and restart services
echo "Rebuilding services..."
docker-compose build

echo "Restarting services..."
docker-compose down
docker-compose up -d

# Show logs
echo ""
echo "Update complete! Showing logs..."
docker-compose logs -f --tail=50
