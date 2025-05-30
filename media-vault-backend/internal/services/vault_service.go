package services

import (
	"context"
	"errors"
	"mime/multipart"

	"media-vault/internal/models"
)

type VaultService struct {
    db *sql.DB
}

func NewVaultService(db *sql.DB) *VaultService {
    return &VaultService{db: db}
}

// Add your VaultService methods here
