package handlers

import (
	"context"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"media-vault/internal/models"
	"media-vault/internal/services"
)

type VaultHandler struct {
    vaultService *services.VaultService
}

func NewVaultHandler(vaultService *services.VaultService) *VaultHandler {
    return &VaultHandler{vaultService: vaultService}
}

// Add your VaultHandler methods here
