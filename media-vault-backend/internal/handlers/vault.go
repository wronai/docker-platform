package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wronai/media-vault-backend/internal/services"
)

type VaultHandler struct {
    vaultService *services.VaultService
}

func NewVaultHandler(vaultService *services.VaultService) *VaultHandler {
    return &VaultHandler{vaultService: vaultService}
}

// GetVault returns the user's vault contents
// @Summary Get vault contents
// @Description Get the authenticated user's vault contents
// @Tags vault
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /vault [get]
func (h *VaultHandler) GetVault(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// TODO: Implement actual vault retrieval logic
	// For now, return a mock response
	return c.JSON(fiber.Map{
		"user_id": userID,
		"items":    []string{},
		"quota":    "0/1024 MB",
	})
}
