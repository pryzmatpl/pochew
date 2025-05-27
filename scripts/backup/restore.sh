#!/bin/bash

set -e

# Configuration
BACKUP_DIR="/app/backups"

echo "Starting restore process at $(date)"

# Function to list available backups
list_backups() {
    echo "Available database backups:"
    ls -la "$BACKUP_DIR"/*.sql 2>/dev/null || echo "No database backups found"
    echo ""
    echo "Available storage backups:"
    ls -la "$BACKUP_DIR"/*.tar.gz 2>/dev/null || echo "No storage backups found"
}

# Function to restore database
restore_database() {
    local backup_file="$1"
    if [ ! -f "$backup_file" ]; then
        echo "Error: Database backup file not found: $backup_file"
        exit 1
    fi
    
    echo "Restoring database from: $backup_file"
    psql "$DATABASE_URL" < "$backup_file"
    echo "Database restore completed"
}

# Function to restore storage
restore_storage() {
    local backup_file="$1"
    if [ ! -f "$backup_file" ]; then
        echo "Error: Storage backup file not found: $backup_file"
        exit 1
    fi
    
    echo "Restoring storage from: $backup_file"
    tar -xzf "$backup_file" -C /app/
    echo "Storage restore completed"
}

# Main restore logic
if [ $# -eq 0 ]; then
    echo "Usage: $0 [database_backup_file] [storage_backup_file]"
    echo "       $0 list  # to list available backups"
    exit 1
fi

if [ "$1" = "list" ]; then
    list_backups
    exit 0
fi

# Restore database if provided
if [ -n "$1" ]; then
    restore_database "$BACKUP_DIR/$1"
fi

# Restore storage if provided
if [ -n "$2" ]; then
    restore_storage "$BACKUP_DIR/$2"
fi

echo "Restore process completed at $(date)" 