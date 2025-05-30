package services

import (
	"database/sql"
)

type VaultService struct {
    db *sql.DB
}

func NewVaultService(db *sql.DB) *VaultService {
    return &VaultService{db: db}
}

// Add your VaultService methods here
