package handlers

import (
	"context"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"media-vault-backend/internal/models"
	"media-vault-backend/internal/services"
)

type VaultHandler struct {
    vaultService *services.VaultService
}

func NewVaultHandler(vaultService *services.VaultService) *VaultHandler {
    return &VaultHandler{vaultService: vaultService}
}

// Add your VaultHandler methods here
