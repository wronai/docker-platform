#!/bin/bash

# Fix user.go
cat > media-vault-backend/internal/models/user.go << 'EOL'
package models

import "time"

// User represents a user in the system
type User struct {
    ID        string    `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Username  string    `json:"username" db:"username"`
    FullName  string    `json:"full_name,omitempty" db:"full_name"`
    AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"`
    IsAdmin   bool      `json:"is_admin" db:"is_admin"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserSettings represents user preferences and settings
type UserSettings struct {
    UserID         string    `json:"user_id" db:"user_id"`
    Theme          string    `json:"theme" db:"theme"`
    Notifications  bool      `json:"notifications" db:"notifications"`
    EmailAlerts    bool      `json:"email_alerts" db:"email_alerts"`
    StorageQuota   int64     `json:"storage_quota" db:"storage_quota"`
    StorageUsed    int64     `json:"storage_used" db:"storage_used"`
    LastBackup     *time.Time `json:"last_backup,omitempty" db:"last_backup"`
    AutoSave       bool      `json:"auto_save" db:"auto_save"`
    AutoTagging    bool      `json:"auto_tagging" db:"auto_tagging"`
    AutoCategorize bool      `json:"auto_categorize" db:"auto_categorize"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
EOL

# Fix vault_service.go
cat > media-vault-backend/internal/services/vault_service.go << 'EOL'
package services

import (
	"database/sql"
	"media-vault-backend/internal/models"
)

type VaultService struct {
    db *sql.DB
}

func NewVaultService(db *sql.DB) *VaultService {
    return &VaultService{db: db}
}

// Add your VaultService methods here
EOL

# Fix vault.go
cat > media-vault-backend/internal/handlers/vault.go << 'EOL'
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"media-vault-backend/internal/services"
)

type VaultHandler struct {
    vaultService *services.VaultService
}

func NewVaultHandler(vaultService *services.VaultService) *VaultHandler {
    return &VaultHandler{vaultService: vaultService}
}

// Add your VaultHandler methods here
EOL

echo "Fixed Go files with proper package declarations and basic structure"
