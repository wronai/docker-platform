#!/bin/bash

# Enhanced backup script with S3 support
set -e

BACKUP_DIR="/backups"
SOURCE_DATA="/backup-source/data"
SOURCE_UPLOADS="/backup-source/uploads"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="media-vault-backup-${TIMESTAMP}"

echo "🗄️  Starting backup: $BACKUP_NAME"

# Create backup directory
mkdir -p "$BACKUP_DIR/$BACKUP_NAME"

# Backup SQLite database
echo "📄 Backing up database..."
if [ -f "$SOURCE_DATA/media.db" ]; then
    sqlite3 "$SOURCE_DATA/media.db" ".backup $BACKUP_DIR/$BACKUP_NAME/media.db"
    echo "✅ Database backup completed"
else
    echo "⚠️  Database not found"
fi

# Backup uploads
echo "📁 Backing up uploads..."
if [ -d "$SOURCE_UPLOADS" ]; then
    cp -r "$SOURCE_UPLOADS" "$BACKUP_DIR/$BACKUP_NAME/"
    echo "✅ Uploads backup completed"
else
    echo "⚠️  Uploads directory not found"
fi

# Create metadata
cat > "$BACKUP_DIR/$BACKUP_NAME/metadata.json" << EOF
{
  "backup_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "version": "1.0",
  "components": ["database", "uploads"],
  "size_bytes": $(du -sb "$BACKUP_DIR/$BACKUP_NAME" | cut -f1)
}
EOF

# Compress backup
echo "🗜️  Compressing backup..."
cd "$BACKUP_DIR"
tar -czf "${BACKUP_NAME}.tar.gz" "$BACKUP_NAME"
rm -rf "$BACKUP_NAME"

# Upload to S3 if configured
if [ -n "$S3_BUCKET" ] && [ -n "$AWS_ACCESS_KEY_ID" ]; then
    echo "☁️  Uploading to S3..."
    aws s3 cp "${BACKUP_NAME}.tar.gz" "s3://$S3_BUCKET/media-vault-backups/"
    echo "✅ S3 upload completed"
fi

# Cleanup old backups
if [ -n "$RETENTION_DAYS" ]; then
    echo "🧹 Cleaning up old backups..."
    find "$BACKUP_DIR" -name "media-vault-backup-*.tar.gz" -mtime +$RETENTION_DAYS -delete
    echo "✅ Cleanup completed"
fi

echo "🎉 Backup completed: ${BACKUP_NAME}.tar.gz"

