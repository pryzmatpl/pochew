#!/bin/bash

set -e

# Configuration
BACKUP_DIR="/app/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_BACKUP_FILE="db_backup_${TIMESTAMP}.sql"
STORAGE_BACKUP_FILE="storage_backup_${TIMESTAMP}.tar.gz"

# Ensure backup directory exists
mkdir -p "$BACKUP_DIR"

echo "Starting backup process at $(date)"

# Database backup
echo "Creating database backup..."
pg_dump "$DATABASE_URL" > "$BACKUP_DIR/$DB_BACKUP_FILE"
echo "Database backup completed: $DB_BACKUP_FILE"

# Storage backup
if [ -d "/app/storage" ]; then
    echo "Creating storage backup..."
    tar -czf "$BACKUP_DIR/$STORAGE_BACKUP_FILE" -C /app storage/
    echo "Storage backup completed: $STORAGE_BACKUP_FILE"
fi

# Cleanup old backups based on retention policy
if [ -n "$BACKUP_RETENTION" ] && [ "$BACKUP_RETENTION" -gt 0 ]; then
    echo "Cleaning up old backups (keeping last $BACKUP_RETENTION days)..."
    find "$BACKUP_DIR" -name "*.sql" -mtime +$BACKUP_RETENTION -delete
    find "$BACKUP_DIR" -name "*.tar.gz" -mtime +$BACKUP_RETENTION -delete
fi

echo "Backup process completed at $(date)"

# If running as a scheduled service, sleep until next run
if [ -n "$BACKUP_SCHEDULE" ]; then
    echo "Backup service running in scheduled mode"
    while true; do
        sleep 3600  # Sleep for 1 hour, cron-like scheduling would be handled externally
    done
fi 