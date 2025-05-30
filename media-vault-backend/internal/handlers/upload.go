package handlers

import (
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/wronai/media-vault-backend/internal/models"
	"github.com/wronai/media-vault-backend/internal/services"
)

// UploadHandler handles file uploads
// @title Media Vault Upload Handler
// @version 1.0
// @description Handles file uploads to the media vault
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type UploadHandler struct {
	vaultService       services.VaultService
	photoService       services.PhotoService
	descriptionService services.DescriptionService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(
	vaultService services.VaultService,
	photoService services.PhotoService,
	descriptionService services.DescriptionService,
) *UploadHandler {
	return &UploadHandler{
		vaultService:       vaultService,
		photoService:       photoService,
		descriptionService: descriptionService,
	}
}

// UploadSingle handles single file upload
// @Summary Upload a single file
// @Description Upload a single file to the media vault
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param description formData string false "File description"
// @Param tags formData string false "Comma-separated list of tags"
// @Success 200 {object} models.Photo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vault/upload [post]
func (h *UploadHandler) UploadSingle(c *fiber.Ctx) error {
	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Check if the request contains a file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded: " + err.Error(),
		})
	}

	// Get additional form data
	description := c.FormValue("description", "")
	tags := c.FormValue("tags", "")

	// Check file size (max 10MB)
	if fileHeader.Size > 10<<20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large. Maximum size is 10MB",
		})
	}

	// Check file type
	ext := filepath.Ext(fileHeader.Filename)
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedTypes[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Only JPG, JPEG, PNG, GIF, and WebP are allowed",
		})
	}

	// Prepare metadata
	metadata := make(map[string]interface{})
	if description != "" {
		metadata["description"] = description
	}
	if tags != "" {
		metadata["tags"] = tags
	}

	// Upload the file
	photo, err := h.photoService.UploadPhoto(c.Context(), userID, fileHeader, metadata)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(photo)
}

// BulkUpload handles multiple file uploads
// @Summary Upload multiple files
// @Description Upload multiple files to the media vault
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Files to upload"
// @Param description formData string false "Description for all files"
// @Param tags formData string false "Comma-separated list of tags for all files"
// @Success 200 {array} models.Photo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vault/upload/bulk [post]
func (h *UploadHandler) BulkUpload(c *fiber.Ctx) error {
	// Get user ID from context
	userID := c.Locals("userID").(string)

	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form: " + err.Error(),
		})
	}

	// Get the files
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	// Get additional form data
	description := c.FormValue("description", "")
	tags := c.FormValue("tags", "")

	// Process each file
	var uploadedPhotos []*models.Photo
	for _, fileHeader := range files {
		// Check file size (max 10MB)
		if fileHeader.Size > 10<<20 {
			continue // Skip files that are too large
		}

		// Check file type
		ext := filepath.Ext(fileHeader.Filename)
		allowedTypes := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
		}

		if !allowedTypes[ext] {
			continue // Skip unsupported file types
		}

		// Upload the file
		photo, err := h.photoService.UploadPhoto(c.Context(), fileHeader, userID, "")
		if err != nil {
			log.Printf("Failed to upload file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Update photo with additional metadata if provided
		if description != "" || tags != "" {
			updateData := map[string]interface{}{
				"description": description,
			}

			if tags != "" {
				updateData["tags"] = tags
			}

			_, err = h.photoService.UpdatePhoto(c.Context(), photo.ID, updateData)
			if err != nil {
				log.Printf("Failed to update photo metadata: %v", err)
			}

			// Update the photo object with the new data
			if description != "" {
				desc := description
				photo.Description = &desc
			}
			if tags != "" {
				tagsStr := tags
				photo.Tags = &tagsStr
			}
		}

		uploadedPhotos = append(uploadedPhotos, photo)
	}

	if len(uploadedPhotos) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files were uploaded. Make sure files are valid images and under 10MB",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(uploadedPhotos)
}

// Helper function to get content type from file extension
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}
	return contentTypes[ext]
}

// Helper function to check if a file is an image
func isImage(contentType string) bool {
	return contentType == "image/jpeg" || 
	       contentType == "image/png" || 
	       contentType == "image/gif"
}

// Helper function to validate image file
func validateImageFile(fileHeader *multipart.FileHeader) error {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the first 512 bytes to detect the content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Reset the file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Detect content type
	contentType := http.DetectContentType(buffer)
	if !isImage(contentType) {
		return fmt.Errorf("invalid file type: %s", contentType)
	}

	return nil
}