package handlers

import (
	"github.com/gofiber/fiber/v2"
	"media-vault/internal/services"
)

// PhotoHandler handles photo-related HTTP requests
type PhotoHandler struct {
	photoService *services.PhotoService
}

// NewPhotoHandler creates a new PhotoHandler
func NewPhotoHandler(photoService *services.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

// UploadPhoto handles photo uploads
func (h *PhotoHandler) UploadPhoto(c *fiber.Ctx) error {
	// Get the file from the form data
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID").(string)

	// Upload the photo
	photo, err := h.photoService.UploadPhoto(c.Context(), userID, file, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload photo: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(photo)
}

// GetPhoto handles retrieving a photo by ID
func (h *PhotoHandler) GetPhoto(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Get photo from service
	photo, err := h.photoService.GetPhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	// Check if user has permission to view this photo
	// (Implementation depends on your auth system)


	return c.JSON(photo)
}

// ListPhotos handles listing photos with pagination
func (h *PhotoHandler) ListPhotos(c *fiber.Ctx) error {
	// Get pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Get photos from service
	photos, total, err := h.photoService.ListPhotos(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch photos",
		})
	}

	return c.JSON(fiber.Map{
		"data":  photos,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// DeletePhoto handles deleting a photo
func (h *PhotoHandler) DeletePhoto(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Delete the photo
	err := h.photoService.DeletePhoto(c.Context(), userID, photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete photo: " + err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RegisterRoutes registers photo-related routes
func (h *PhotoHandler) RegisterRoutes(router fiber.Router) {
	photos := router.Group("/photos")
	photos.Post("/", h.UploadPhoto)
	photos.Get("/", h.ListPhotos)
	photos.Get("/:id", h.GetPhoto)
	photos.Delete("/:id", h.DeletePhoto)
}