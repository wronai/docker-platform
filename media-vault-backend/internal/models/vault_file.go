package models

import "time"

// VaultFile represents a file in the vault
type VaultFile struct {
    ID          string    `json:"id" db:"id"`
    UserID      string    `json:"user_id" db:"user_id"`
    Name        string    `json:"name" db:"name"`
    Path        string    `json:"path" db:"path"`
    Size        int64     `json:"size" db:"size"`
    MimeType    string    `json:"mime_type" db:"mime_type"`
    IsPublic    bool      `json:"is_public" db:"is_public"`
    IsEncrypted bool      `json:"is_encrypted" db:"is_encrypted"`
    Metadata    string    `json:"metadata,omitempty" db:"metadata"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
