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
			"error": "Failed to fetch photos: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  photos,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdatePhoto handles updating a photo's metadata
func (h *PhotoHandler) UpdatePhoto(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Parse request body
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Update the photo
	photo, err := h.photoService.UpdatePhoto(c.Context(), userID, photoID, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update photo: " + err.Error(),
		})
	}

	return c.JSON(photo)
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

// GetThumbnail handles retrieving a photo's thumbnail
func (h *PhotoHandler) GetThumbnail(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Get the thumbnail data
	thumbnail, err := h.photoService.GetThumbnail(c.Context(), userID, photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get thumbnail: " + err.Error(),
		})
	}

	// Set appropriate content type
	c.Set("Content-Type", "image/jpeg") // or determine from thumbnail data
	return c.Send(thumbnail)
}

// UpdateDescription updates a photo's description
func (h *PhotoHandler) UpdateDescription(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Parse request body
	var request struct {
		Description string `json:"description"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Update the photo description
	updateData := map[string]interface{}{
		"description": request.Description,
	}

	photo, err := h.photoService.UpdatePhoto(c.Context(), userID, photoID, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update description: " + err.Error(),
		})
	}

	return c.JSON(photo)
}

// GenerateDescription generates an AI description for a photo
func (h *PhotoHandler) GenerateDescription(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Generate the description
	description, err := h.photoService.GenerateDescription(c.Context(), userID, photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate description: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"description": description,
	})
}

// GetSharedWith gets the list of users a photo is shared with
func (h *PhotoHandler) GetSharedWith(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Get the list of users the photo is shared with
	sharedWith, err := h.photoService.GetSharedWith(c.Context(), userID, photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get shared with list: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"shared_with": sharedWith,
	})
}

// RegisterRoutes registers photo-related routes
func (h *PhotoHandler) RegisterRoutes(router fiber.Router) {
	photos := router.Group("/photos")
	photos.Post("/", h.UploadPhoto)
	photos.Get("/", h.ListPhotos)
	photos.Get("/:id", h.GetPhoto)
	photos.Put("/:id", h.UpdatePhoto)
	photos.Delete("/:id", h.DeletePhoto)
	photos.Get("/:id/thumbnail", h.GetThumbnail)
	photos.Put("/:id/description", h.UpdateDescription)
	photos.Post("/:id/generate-description", h.GenerateDescription)
	photos.Get("/:id/shared-with", h.GetSharedWith)
}