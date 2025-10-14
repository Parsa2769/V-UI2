#!/bin/bash

set -e

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$TIMESTAMP.tar.gz"

echo "Creating backup..."

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Export database
echo "Exporting database..."
docker-compose exec -T postgres pg_dump -U x3ui x3ui > "$BACKUP_DIR/db_$TIMESTAMP.sql"

# Backup data directory and configs
echo "Backing up data and configs..."
tar -czf $BACKUP_FILE \
    ./data \
    ./config \
    .env \
    "$BACKUP_DIR/db_$TIMESTAMP.sql"

# Remove temporary SQL dump
rm "$BACKUP_DIR/db_$TIMESTAMP.sql"

# Clean old backups (keep last 30 days)
find $BACKUP_DIR -name "backup_*.tar.gz" -mtime +30 -delete

echo "✓ Backup created: $BACKUP_FILE"
echo "  Size: $(du -h $BACKUP_FILE | cut -f1)"
