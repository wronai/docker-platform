package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wronai/media-vault-backend/internal/services"
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
			"error": "Photo not found: " + err.Error(),
		})
	}

	// Check if user has permission to view this photo
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to view this photo",
		})
	}

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
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Update photo
	photo, err := h.photoService.UpdatePhoto(c.Context(), photoID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update photo: " + err.Error(),
		})
	}

	// Verify ownership
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this photo",
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

	// First get the photo to verify ownership
	photo, err := h.photoService.GetPhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found: " + err.Error(),
		})
	}

	// Verify ownership
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this photo",
		})
	}

	// Delete the photo
	err = h.photoService.DeletePhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete photo: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetThumbnail handles retrieving a photo's thumbnail
func (h *PhotoHandler) GetThumbnail(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// First get the photo to verify ownership
	photo, err := h.photoService.GetPhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found: " + err.Error(),
		})
	}

	// Verify ownership
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to view this thumbnail",
		})
	}

	// Get size parameter (optional, default to medium)
	size := c.Query("size", "medium")

	// Get thumbnail data
	data, contentType, err := h.photoService.GetThumbnail(c.Context(), photoID, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get thumbnail: " + err.Error(),
		})
	}

	// Set content type and send the image data
	c.Set("Content-Type", contentType)
	return c.Send(data)
}

// UpdateDescription updates a photo's description
func (h *PhotoHandler) UpdateDescription(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// First get the photo to verify ownership
	photo, err := h.photoService.GetPhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found: " + err.Error(),
		})
	}

	// Verify ownership
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this photo's description",
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

	// Update the description
	updates := map[string]interface{}{
		"description": request.Description,
	}

	// Update the photo
	updatedPhoto, err := h.photoService.UpdatePhoto(c.Context(), photoID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update description: " + err.Error(),
		})
	}

	return c.JSON(updatedPhoto)
}

// GenerateDescription generates an AI description for a photo
func (h *PhotoHandler) GenerateDescription(c *fiber.Ctx) error {
	photoID := c.Params("id")
	if photoID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo ID is required",
		})
	}

	// First get the photo to verify ownership
	photo, err := h.photoService.GetPhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found: " + err.Error(),
		})
	}

	// Verify ownership
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to generate a description for this photo",
		})
	}

	// Generate the description
	description, err := h.photoService.GenerateDescription(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate description: " + err.Error(),
		})
	}

	// Update the photo with the generated description
	updates := map[string]interface{}{
		"description": description,
	}

	_, err = h.photoService.UpdatePhoto(c.Context(), photoID, updates)
	if err != nil {
		// Log the error but don't fail the request
		// since we still want to return the generated description
		// Consider using a background job for this in production
		// log.Printf("Failed to update photo with generated description: %v", err)
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

	// First get the photo to verify ownership
	photo, err := h.photoService.GetPhoto(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found: " + err.Error(),
		})
	}

	// Verify ownership
	userID := c.Locals("userID").(string)
	if photo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to view sharing information for this photo",
		})
	}

	// Get the list of users the photo is shared with
	users, err := h.photoService.GetSharedWith(c.Context(), photoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get shared users: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"photo_id":   photoID,
		"shared_with": users,
	})
}

// RegisterRoutes registers photo-related routes
func (h *PhotoHandler) RegisterRoutes(router fiber.Router) {
	photoGroup := router.Group("/photos")
	{
		// Upload a new photo
		photoGroup.Post("/", h.UploadPhoto)
		
		// List all photos with pagination
		photoGroup.Get("/", h.ListPhotos)
		
		// Get a specific photo by ID
		photoGroup.Get("/:id", h.GetPhoto)
		
		// Update a photo's metadata
		photoGroup.Put("/:id", h.UpdatePhoto)
		
		// Delete a photo
		photoGroup.Delete("/:id", h.DeletePhoto)
		
		// Get a photo's thumbnail
		photoGroup.Get("/:id/thumbnail", h.GetThumbnail)
		
		// Update a photo's description
		photoGroup.Put("/:id/description", h.UpdateDescription)
		
		// Generate an AI description for a photo
		photoGroup.Post("/:id/generate-description", h.GenerateDescription)
		
		// Get the list of users a photo is shared with
		photoGroup.Get("/:id/shared-with", h.GetSharedWith)
	}
}