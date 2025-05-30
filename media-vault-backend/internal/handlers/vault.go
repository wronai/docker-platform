package handlers

import (
	"github.com/wronai/media-vault-backend/internal/services"
)

type VaultHandler struct {
    vaultService *services.VaultService
}

func NewVaultHandler(vaultService *services.VaultService) *VaultHandler {
    return &VaultHandler{vaultService: vaultService}
}

// Add your VaultHandler methods here
