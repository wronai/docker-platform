package database

import (
    "database/sql"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

func Initialize() (*sql.DB, error) {
    dbPath := os.Getenv("DATABASE_PATH")
    if dbPath == "" {
        dbPath = "./data/media.db"
    }

    // Ensure directory exists
    dir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // Enable foreign keys
    if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
        return nil, err
    }

    // Run migrations
    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return db, nil
}

func runMigrations(db *sql.DB) error {
    migrations := []string{
        `CREATE TABLE IF NOT EXISTS photos (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            partner_id TEXT,
            filename TEXT NOT NULL,
            original_name TEXT NOT NULL,
            file_path TEXT NOT NULL,
            thumbnail_path TEXT,
            file_size INTEGER NOT NULL,
            mime_type TEXT NOT NULL,
            width INTEGER,
            height INTEGER,
            hash TEXT NOT NULL,
            description TEXT,
            ai_description TEXT,
            tags TEXT,
            ai_confidence REAL,
            is_nsfw BOOLEAN,
            nsfw_confidence REAL,
            moderation_status TEXT DEFAULT 'pending',
            exif_data TEXT,
            location TEXT,
            camera_make TEXT,
            camera_model TEXT,
            taken_at DATETIME,
            is_shared BOOLEAN DEFAULT FALSE,
            share_count INTEGER DEFAULT 0,
            view_count INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            processed_at DATETIME
        )`,

        `CREATE TABLE IF NOT EXISTS photo_sharing (
            id TEXT PRIMARY KEY,
            photo_id TEXT NOT NULL,
            shared_by TEXT NOT NULL,
            shared_with TEXT NOT NULL,
            permission TEXT NOT NULL DEFAULT 'view',
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            expires_at DATETIME,
            FOREIGN KEY (photo_id) REFERENCES photos (id) ON DELETE CASCADE
        )`,

        `CREATE TABLE IF NOT EXISTS photo_analytics (
            photo_id TEXT PRIMARY KEY,
            views INTEGER DEFAULT 0,
            downloads INTEGER DEFAULT 0,
            shares INTEGER DEFAULT 0,
            last_viewed DATETIME,
            FOREIGN KEY (photo_id) REFERENCES photos (id) ON DELETE CASCADE
        )`,

        `CREATE INDEX IF NOT EXISTS idx_photos_user_id ON photos (user_id)`,
        `CREATE INDEX IF NOT EXISTS idx_photos_partner_id ON photos (partner_id)`,
        `CREATE INDEX IF NOT EXISTS idx_photos_created_at ON photos (created_at)`,
        `CREATE INDEX IF NOT EXISTS idx_photos_moderation_status ON photos (moderation_status)`,
        `CREATE INDEX IF NOT EXISTS idx_photo_sharing_photo_id ON photo_sharing (photo_id)`,
        `CREATE INDEX IF NOT EXISTS idx_photo_sharing_shared_with ON photo_sharing (shared_with)`,
    }

    for _, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return err
        }
    }

    return nil
}

