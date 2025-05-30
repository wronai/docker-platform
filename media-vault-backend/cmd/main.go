// github.com/wronai/media-vault-backend/cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/wronai/media-vault-backend/internal/auth"
	"github.com/wronai/media-vault-backend/internal/database"
	"github.com/wronai/media-vault-backend/internal/handlers"
	"github.com/wronai/media-vault-backend/internal/services"
	"github.com/wronai/media-vault-backend/internal/models"
)

func main() {
	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Keycloak authentication
	if os.Getenv("OAUTH2_ENABLED") == "true" {
		if err := auth.InitKeycloak(); err != nil {
			log.Printf("Warning: Failed to initialize Keycloak: %v", err)
			log.Println("Running without authentication")
		} else {
			log.Println("Keycloak authentication enabled")
		}
	}

	// Initialize services
	vaultService := services.NewVaultService(db)
	photoService := services.NewPhotoService(db)
	descriptionService := services.NewDescriptionService()
	sharingService := services.NewSharingService(db)

	// Initialize handlers
	vaultHandler := handlers.NewVaultHandler(vaultService, photoService)
	adminHandler := handlers.NewAdminHandler(vaultService)
	partnerHandler := handlers.NewPartnerHandler(photoService, sharingService)
	uploadHandler := handlers.NewUploadHandler(vaultService, photoService, descriptionService)
	photoHandler := handlers.NewPhotoHandler(photoService)

	// Fiber app configuration
	app := fiber.New(fiber.Config{
		BodyLimit:    100 * 1024 * 1024, // 100MB
		ReadTimeout:  0,
		WriteTimeout: 0,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
				"code":  code,
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ORIGINS"),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
	}))

	// Health and monitoring endpoints
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "media-vault-api",
			"version": os.Getenv("VAULT_VERSION"),
		})
	})
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Media Vault Metrics"}))

	// Public routes
	app.Get("/api/vault/info", vaultHandler.GetVaultInfo)

	// Protected routes
	api := app.Group("/api")
	api.Use(auth.JWTMiddleware())

	// User endpoints
	api.Get("/user", auth.GetUserInfo)
	api.Get("/user/profile", vaultHandler.GetUserProfile)
	api.Put("/user/profile", vaultHandler.UpdateUserProfile)

	// Vault endpoints (User level)
	vault := api.Group("/vault")
	vault.Post("/upload", uploadHandler.UploadSingle)
	vault.Post("/upload/bulk", uploadHandler.BulkUpload)
	vault.Get("/files", vaultHandler.GetUserFiles)
	vault.Get("/files/:id", vaultHandler.GetFile)
	vault.Put("/files/:id", vaultHandler.UpdateFile)
	vault.Delete("/files/:id", vaultHandler.DeleteFile)
	vault.Get("/files/:id/download", vaultHandler.DownloadFile)
	vault.Get("/files/:id/thumbnail", vaultHandler.GetThumbnail)
	vault.Get("/search", vaultHandler.SearchFiles)
	vault.Get("/stats", vaultHandler.GetUserStats)

	// Photo management endpoints
	photos := api.Group("/photos")
	photos.Get("/", photoHandler.ListPhotos)
	photos.Get("/:id", photoHandler.GetPhoto)
	photos.Put("/:id", photoHandler.UpdatePhoto)
	photos.Delete("/:id", photoHandler.DeletePhoto)
	photos.Get("/:id/thumbnail", photoHandler.GetThumbnail)
	photos.Put("/:id/description", photoHandler.UpdateDescription)
	photos.Post("/:id/generate-description", photoHandler.GenerateDescription)
	photos.Get("/:id/shared-with", photoHandler.GetSharedWith)

	// Partner endpoints (requires partner role)
	partner := api.Group("/partner")
	partner.Use(auth.RequireRole("vault-partner"))
	partner.Post("/bulk-upload", partnerHandler.BulkUpload)
	partner.Get("/photos", partnerHandler.GetPartnerPhotos)
	partner.Put("/photos/batch-update", partnerHandler.BatchUpdateDescriptions)
	partner.Post("/photos/batch-share", partnerHandler.BatchSharePhotos)
	partner.Get("/photos/:id/analytics", partnerHandler.GetPhotoAnalytics)
	partner.Get("/dashboard", partnerHandler.GetPartnerDashboard)
	partner.Get("/analytics", partnerHandler.GetPartnerAnalytics)

	// Admin endpoints (requires admin role)
	admin := api.Group("/admin")
	admin.Use(auth.RequireRole("vault-admin"))
	admin.Get("/users", adminHandler.GetUsers)
	admin.Get("/users/:id", adminHandler.GetUser)
	admin.Put("/users/:id", adminHandler.UpdateUser)
	admin.Delete("/users/:id", adminHandler.DeleteUser)
	admin.Get("/users/:id/files", adminHandler.GetUserFiles)
	admin.Get("/system/stats", adminHandler.GetSystemStats)
	admin.Get("/system/health", adminHandler.GetSystemHealth)
	admin.Get("/content/flagged", adminHandler.GetFlaggedContent)
	admin.Put("/content/:id/moderate", adminHandler.ModerateContent)
	admin.Get("/audit-logs", adminHandler.GetAuditLogs)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Media Vault API starting on port %s", port)
	log.Printf("Vault Name: %s", os.Getenv("VAULT_NAME"))
	log.Printf("Environment: %s", os.Getenv("ENVIRONMENT"))

	log.Fatal(app.Listen(":" + port))
}

---

// github.com/wronai/media-vault-backend/internal/models/photo.go
package models

import (
	"time"
)

// Photo represents a photo in the vault
type Photo struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	PartnerID     *string   `json:"partner_id,omitempty" db:"partner_id"`
	Filename      string    `json:"filename" db:"filename"`
	OriginalName  string    `json:"original_name" db:"original_name"`
	FilePath      string    `json:"file_path" db:"file_path"`
	ThumbnailPath *string   `json:"thumbnail_path,omitempty" db:"thumbnail_path"`
	FileSize      int64     `json:"file_size" db:"file_size"`
	MimeType      string    `json:"mime_type" db:"mime_type"`
	Width         *int      `json:"width,omitempty" db:"width"`
	Height        *int      `json:"height,omitempty" db:"height"`
	Hash          string    `json:"hash" db:"hash"`

	// AI Generated Content
	Description         *string `json:"description,omitempty" db:"description"`
	AIDescription       *string `json:"ai_description,omitempty" db:"ai_description"`
	Tags               *string `json:"tags,omitempty" db:"tags"`
	AIConfidence       *float64 `json:"ai_confidence,omitempty" db:"ai_confidence"`

	// Content Moderation
	IsNSFW             *bool    `json:"is_nsfw,omitempty" db:"is_nsfw"`
	NSFWConfidence     *float64 `json:"nsfw_confidence,omitempty" db:"nsfw_confidence"`
	ModerationStatus   string   `json:"moderation_status" db:"moderation_status"` // pending, approved, rejected

	// Metadata
	ExifData           *string `json:"exif_data,omitempty" db:"exif_data"`
	Location           *string `json:"location,omitempty" db:"location"`
	CameraMake         *string `json:"camera_make,omitempty" db:"camera_make"`
	CameraModel        *string `json:"camera_model,omitempty" db:"camera_model"`
	TakenAt            *time.Time `json:"taken_at,omitempty" db:"taken_at"`

	// Sharing
	IsShared           bool     `json:"is_shared" db:"is_shared"`
	SharedWith         []string `json:"shared_with,omitempty"`
	ShareCount         int      `json:"share_count" db:"share_count"`
	ViewCount          int      `json:"view_count" db:"view_count"`

	// Timestamps
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	ProcessedAt        *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// PhotoSharing represents photo sharing permissions
type PhotoSharing struct {
	ID        string    `json:"id" db:"id"`
	PhotoID   string    `json:"photo_id" db:"photo_id"`
	SharedBy  string    `json:"shared_by" db:"shared_by"`
	SharedWith string   `json:"shared_with" db:"shared_with"`
	Permission string   `json:"permission" db:"permission"` // view, download
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// PhotoAnalytics represents photo analytics data
type PhotoAnalytics struct {
	PhotoID     string    `json:"photo_id" db:"photo_id"`
	Views       int       `json:"views" db:"views"`
	Downloads   int       `json:"downloads" db:"downloads"`
	Shares      int       `json:"shares" db:"shares"`
	LastViewed  *time.Time `json:"last_viewed,omitempty" db:"last_viewed"`
	ViewerIPs   []string  `json:"viewer_ips,omitempty"`
	Countries   []string  `json:"countries,omitempty"`
}

// BulkUploadRequest represents a bulk upload request
type BulkUploadRequest struct {
	Files       []UploadedFile `json:"files"`
	DefaultTags string         `json:"default_tags,omitempty"`
	ShareWith   []string       `json:"share_with,omitempty"`
	AutoDescribe bool          `json:"auto_describe"`
}

// UploadedFile represents a file in bulk upload
type UploadedFile struct {
	Filename    string `json:"filename"`
	Data        []byte `json:"data"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
}

// BatchUpdateRequest represents batch description update
type BatchUpdateRequest struct {
	PhotoIDs    []string `json:"photo_ids"`
	Description string   `json:"description,omitempty"`
	Tags        string   `json:"tags,omitempty"`
	Operation   string   `json:"operation"` // replace, append, prepend
}

// ShareRequest represents a photo sharing request
type ShareRequest struct {
	PhotoIDs   []string   `json:"photo_ids"`
	ShareWith  []string   `json:"share_with"`
	Permission string     `json:"permission"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Message    string     `json:"message,omitempty"`
}

---

// github.com/wronai/media-vault-backend/internal/handlers/partner.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wronai/media-vault-backend/internal/auth"
	"github.com/wronai/media-vault-backend/internal/models"
	"github.com/wronai/media-vault-backend/internal/services"
)

type PartnerHandler struct {
	photoService   *services.PhotoService
	sharingService *services.SharingService
}

func NewPartnerHandler(photoService *services.PhotoService, sharingService *services.SharingService) *PartnerHandler {
	return &PartnerHandler{
		photoService:   photoService,
		sharingService: sharingService,
	}
}

// BulkUpload handles bulk photo upload for partners
func (h *PartnerHandler) BulkUpload(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse form data"})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No files provided"})
	}

	// Get optional parameters
	defaultTags := c.FormValue("default_tags", "")
	autoDescribe := c.FormValue("auto_describe", "false") == "true"
	shareWithStr := c.FormValue("share_with", "")

	var shareWith []string
	if shareWithStr != "" {
		shareWith = strings.Split(shareWithStr, ",")
	}

	var uploadResults []map[string]interface{}
	var successCount, errorCount int

	// Process each file
	for i, file := range files {
		log.Printf("Processing file %d/%d: %s", i+1, len(files), file.Filename)

		// Upload file
		photo, err := h.photoService.CreatePhotoFromUpload(file, user.Subject, &user.Subject)
		if err != nil {
			errorCount++
			uploadResults = append(uploadResults, map[string]interface{}{
				"filename": file.Filename,
				"status":   "error",
				"error":    err.Error(),
			})
			continue
		}

		// Set default tags if provided
		if defaultTags != "" {
			photo.Tags = &defaultTags
		}

		// Generate AI description if requested
		if autoDescribe {
			if description, err := h.photoService.GenerateAIDescription(photo.ID); err == nil {
				photo.AIDescription = &description
			}
		}

		// Update photo with additional data
		if err := h.photoService.UpdatePhoto(photo); err != nil {
			log.Printf("Failed to update photo %s: %v", photo.ID, err)
		}

		// Share with specified users
		if len(shareWith) > 0 {
			shareRequest := &models.ShareRequest{
				PhotoIDs:   []string{photo.ID},
				ShareWith:  shareWith,
				Permission: "view",
			}
			if err := h.sharingService.SharePhotos(user.Subject, shareRequest); err != nil {
				log.Printf("Failed to share photo %s: %v", photo.ID, err)
			}
		}

		successCount++
		uploadResults = append(uploadResults, map[string]interface{}{
			"filename": file.Filename,
			"status":   "success",
			"photo_id": photo.ID,
			"size":     photo.FileSize,
		})
	}

	return c.JSON(fiber.Map{
		"message":       "Bulk upload completed",
		"total_files":   len(files),
		"success_count": successCount,
		"error_count":   errorCount,
		"results":       uploadResults,
	})
}

// GetPartnerPhotos returns photos uploaded by the partner
func (h *PartnerHandler) GetPartnerPhotos(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	sortBy := c.Query("sort", "created_at")
	sortOrder := c.Query("order", "desc")
	search := c.Query("search", "")

	photos, total, err := h.photoService.GetPartnerPhotos(user.Subject, page, limit, sortBy, sortOrder, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get photos"})
	}

	return c.JSON(fiber.Map{
		"photos": photos,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// BatchUpdateDescriptions updates descriptions for multiple photos
func (h *PartnerHandler) BatchUpdateDescriptions(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)

	var request models.BatchUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(request.PhotoIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No photo IDs provided"})
	}

	var updatedCount int
	var errors []string

	for _, photoID := range request.PhotoIDs {
		// Verify partner owns the photo
		photo, err := h.photoService.GetPhoto(photoID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Photo %s not found", photoID))
			continue
		}

		if photo.PartnerID == nil || *photo.PartnerID != user.Subject {
			errors = append(errors, fmt.Sprintf("No permission for photo %s", photoID))
			continue
		}

		// Update description based on operation
		switch request.Operation {
		case "replace":
			photo.Description = &request.Description
			if request.Tags != "" {
				photo.Tags = &request.Tags
			}
		case "append":
			if photo.Description != nil {
				newDesc := *photo.Description + " " + request.Description
				photo.Description = &newDesc
			} else {
				photo.Description = &request.Description
			}
		case "prepend":
			if photo.Description != nil {
				newDesc := request.Description + " " + *photo.Description
				photo.Description = &newDesc
			} else {
				photo.Description = &request.Description
			}
		default:
			photo.Description = &request.Description
		}

		if err := h.photoService.UpdatePhoto(photo); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to update photo %s: %v", photoID, err))
			continue
		}

		updatedCount++
	}

	return c.JSON(fiber.Map{
		"message":       "Batch update completed",
		"updated_count": updatedCount,
		"total_count":   len(request.PhotoIDs),
		"errors":        errors,
	})
}

// BatchSharePhotos shares multiple photos with users
func (h *PartnerHandler) BatchSharePhotos(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)

	var request models.ShareRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(request.PhotoIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No photo IDs provided"})
	}

	if len(request.ShareWith) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "No users to share with"})
	}

	// Verify partner owns all photos
	for _, photoID := range request.PhotoIDs {
		photo, err := h.photoService.GetPhoto(photoID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Photo %s not found", photoID)})
		}

		if photo.PartnerID == nil || *photo.PartnerID != user.Subject {
			return c.Status(403).JSON(fiber.Map{"error": fmt.Sprintf("No permission for photo %s", photoID)})
		}
	}

	// Share photos
	if err := h.sharingService.SharePhotos(user.Subject, &request); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to share photos"})
	}

	return c.JSON(fiber.Map{
		"message":     "Photos shared successfully",
		"photo_count": len(request.PhotoIDs),
		"user_count":  len(request.ShareWith),
	})
}

// GetPhotoAnalytics returns analytics for a specific photo
func (h *PartnerHandler) GetPhotoAnalytics(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)
	photoID := c.Params("id")

	// Verify partner owns the photo
	photo, err := h.photoService.GetPhoto(photoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Photo not found"})
	}

	if photo.PartnerID == nil || *photo.PartnerID != user.Subject {
		return c.Status(403).JSON(fiber.Map{"error": "Permission denied"})
	}

	analytics, err := h.photoService.GetPhotoAnalytics(photoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get analytics"})
	}

	return c.JSON(analytics)
}

// GetPartnerDashboard returns dashboard data for partner
func (h *PartnerHandler) GetPartnerDashboard(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)

	dashboard, err := h.photoService.GetPartnerDashboard(user.Subject)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get dashboard data"})
	}

	return c.JSON(dashboard)
}

// GetPartnerAnalytics returns comprehensive analytics for partner
func (h *PartnerHandler) GetPartnerAnalytics(c *fiber.Ctx) error {
	user := c.Locals("user").(*auth.KeycloakClaims)

	// Parse date range
	fromStr := c.Query("from", "")
	toStr := c.Query("to", "")

	var from, to *time.Time
	if fromStr != "" {
		if parsedFrom, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = &parsedFrom
		}
	}
	if toStr != "" {
		if parsedTo, err := time.Parse("2006-01-02", toStr); err == nil {
			to = &parsedTo
		}
	}

	analytics, err := h.photoService.GetPartnerAnalytics(user.Subject, from, to)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get analytics"})
	}

	return c.JSON(analytics)
}

---

// github.com/wronai/media-vault-backend/internal/services/photo_service.go
package services

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/wronai/media-vault-backend/internal/models"
	"github.com/wronai/media-vault-backend/internal/utils"
)

type PhotoService struct {
	db *sql.DB
}

func NewPhotoService(db *sql.DB) *PhotoService {
	return &PhotoService{db: db}
}

// CreatePhotoFromUpload creates a photo record from uploaded file
func (s *PhotoService) CreatePhotoFromUpload(fileHeader *multipart.FileHeader, userID string, partnerID *string) (*models.Photo, error) {
	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer file.Close()

	// Generate unique filename
	photoID := uuid.New().String()
	ext := filepath.Ext(fileHeader.Filename)
	filename := photoID + ext

	// Create upload directory if not exists
	uploadDir := os.Getenv("UPLOAD_PATH")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	originalDir := filepath.Join(uploadDir, "originals")
	if err := os.MkdirAll(originalDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Save file
	filePath := filepath.Join(originalDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy file and calculate hash
	file.Seek(0, 0) // Reset file pointer
	hasher := md5.New()
	writer := io.MultiWriter(dst, hasher)

	if _, err := io.Copy(writer, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	// Get image dimensions
	file.Seek(0, 0)
	var width, height *int
	if img, _, err := image.Decode(file); err == nil {
		bounds := img.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		width, height = &w, &h
	}

	// Generate thumbnail
	thumbnailPath, err := utils.GenerateThumbnail(filePath, photoID)
	if err != nil {
		// Log error but don't fail upload
		fmt.Printf("Failed to generate thumbnail for %s: %v\n", photoID, err)
	}

	// Extract EXIF data
	exifData, err := utils.ExtractEXIF(filePath)
	if err != nil {
		// Log error but don't fail upload
		fmt.Printf("Failed to extract EXIF for %s: %v\n", photoID, err)
	}

	// Create photo record
	photo := &models.Photo{
		ID:               photoID,
		UserID:           userID,
		PartnerID:        partnerID,
		Filename:         filename,
		OriginalName:     fileHeader.Filename,
		FilePath:         filePath,
		ThumbnailPath:    thumbnailPath,
		FileSize:         fileHeader.Size,
		MimeType:         fileHeader.Header.Get("Content-Type"),
		Width:            width,
		Height:           height,
		Hash:             hash,
		ModerationStatus: "pending",
		IsShared:         false,
		ShareCount:       0,
		ViewCount:        0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ExifData:         exifData,
	}

	// Save to database
	query := `
		INSERT INTO photos (
			id, user_id, partner_id, filename, original_name, file_path,
			thumbnail_path, file_size, mime_type, width, height, hash,
			moderation_status, is_shared, share_count, view_count,
			created_at, updated_at, exif_data
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.Exec(query,
		photo.ID, photo.UserID, photo.PartnerID, photo.Filename,
		photo.OriginalName, photo.FilePath, photo.ThumbnailPath,
		photo.FileSize, photo.MimeType, photo.Width, photo.Height,
		photo.Hash, photo.ModerationStatus, photo.IsShared,
		photo.ShareCount, photo.ViewCount, photo.CreatedAt,
		photo.UpdatedAt, photo.ExifData,
	)

	if err != nil {
		// Clean up uploaded file on database error
		os.Remove(filePath)
		if thumbnailPath != nil {
			os.Remove(*thumbnailPath)
		}
		return nil, fmt.Errorf("failed to save photo to database: %v", err)
	}

	return photo, nil
}

// GetPhoto retrieves a photo by ID
func (s *PhotoService) GetPhoto(photoID string) (*models.Photo, error) {
	query := `
		SELECT id, user_id, partner_id, filename, original_name, file_path,
			   thumbnail_path, file_size, mime_type, width, height, hash,
			   description, ai_description, tags, ai_confidence,
			   is_nsfw, nsfw_confidence, moderation_status,
			   exif_data, location, camera_make, camera_model, taken_at,
			   is_shared, share_count, view_count,
			   created_at, updated_at, processed_at
		FROM photos WHERE id = ?
	`

	photo := &models.Photo{}
	err := s.db.QueryRow(query, photoID).Scan(
		&photo.ID, &photo.UserID, &photo.PartnerID, &photo.Filename,
		&photo.OriginalName, &photo.FilePath, &photo.ThumbnailPath,
		&photo.FileSize, &photo.MimeType, &photo.Width, &photo.Height,
		&photo.Hash, &photo.Description, &photo.AIDescription,
		&photo.Tags, &photo.AIConfidence, &photo.IsNSFW,
		&photo.NSFWConfidence, &photo.ModerationStatus, &photo.ExifData,
		&photo.Location, &photo.CameraMake, &photo.CameraModel,
		&photo.TakenAt, &photo.IsShared, &photo.ShareCount,
		&photo.ViewCount, &photo.CreatedAt, &photo.UpdatedAt,
		&photo.ProcessedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("photo not found")
		}
		return nil, fmt.Errorf("failed to get photo: %v", err)
	}

	return photo, nil
}

// UpdatePhoto updates a photo record
func (s *PhotoService) UpdatePhoto(photo *models.Photo) error {
	query := `
		UPDATE photos SET
			description = ?, ai_description = ?, tags = ?, ai_confidence = ?,
			is_nsfw = ?, nsfw_confidence = ?, moderation_status = ?,
			location = ?, camera_make = ?, camera_model = ?, taken_at = ?,
			is_shared = ?, share_count = ?, view_count = ?,
			updated_at = ?, processed_at = ?
		WHERE id = ?
	`

	photo.UpdatedAt = time.Now()

	_, err := s.db.Exec(query,
		photo.Description, photo.AIDescription, photo.Tags, photo.AIConfidence,
		photo.IsNSFW, photo.NSFWConfidence, photo.ModerationStatus,
		photo.Location, photo.CameraMake, photo.CameraModel, photo.TakenAt,
		photo.IsShared, photo.ShareCount, photo.ViewCount,
		photo.UpdatedAt, photo.ProcessedAt, photo.ID,
	)

	return err
}

// GetPartnerPhotos retrieves photos for a partner with pagination
func (s *PhotoService) GetPartnerPhotos(partnerID string, page, limit int, sortBy, sortOrder, search string) ([]*models.Photo, int, error) {
	offset := (page - 1) * limit

	// Build WHERE clause
	whereClause := "WHERE partner_id = ?"
	args := []interface{}{partnerID}

	if search != "" {
		whereClause += " AND (original_name LIKE ? OR description LIKE ? OR tags LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	// Validate sort parameters
	allowedSorts := map[string]bool{
		"created_at": true, "updated_at": true, "original_name": true,
		"file_size": true, "view_count": true, "share_count": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM photos " + whereClause
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count photos: %v", err)
	}

	// Get photos
	query := fmt.Sprintf(`
		SELECT id, user_id, partner_id, filename, original_name, file_path,
			   thumbnail_path, file_size, mime_type, width, height,
			   description, ai_description, tags, is_nsfw, moderation_status,
			   is_shared, share_count, view_count, created_at, updated_at
		FROM photos %s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, sortBy, sortOrder)

	args = append(args, limit, offset)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query photos: %v", err)
	}
	defer rows.Close()

	var photos []*models.Photo
	for rows.Next() {
		photo := &models.Photo{}
		err := rows.Scan(
			&photo.ID, &photo.UserID, &photo.PartnerID, &photo.Filename,
			&photo.OriginalName, &photo.FilePath, &photo.ThumbnailPath,
			&photo.FileSize, &photo.MimeType, &photo.Width, &photo.Height,
			&photo.Description, &photo.AIDescription, &photo.Tags,
			&photo.IsNSFW, &photo.ModerationStatus, &photo.IsShared,
			&photo.ShareCount, &photo.ViewCount, &photo.CreatedAt, &photo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan photo: %v", err)
		}
		photos = append(photos, photo)
	}

	return photos, total, nil
}

// GenerateAIDescription generates AI description for a photo
func (s *PhotoService) GenerateAIDescription(photoID string) (string, error) {
	// This would integrate with AI service (OpenAI, Google Vision, etc.)
	// For now, return a placeholder
	return "AI-generated description for photo " + photoID, nil
}

// GetPhotoAnalytics retrieves analytics for a photo
func (s *PhotoService) GetPhotoAnalytics(photoID string) (*models.PhotoAnalytics, error) {
	query := `
		SELECT photo_id, views, downloads, shares, last_viewed
		FROM photo_analytics WHERE photo_id = ?
	`

	analytics := &models.PhotoAnalytics{}
	err := s.db.QueryRow(query, photoID).Scan(
		&analytics.PhotoID, &analytics.Views, &analytics.Downloads,
		&analytics.Shares, &analytics.LastViewed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Create default analytics if not exists
			analytics = &models.PhotoAnalytics{
				PhotoID:   photoID,
				Views:     0,
				Downloads: 0,
				Shares:    0,
			}
		} else {
			return nil, fmt.Errorf("failed to get analytics: %v", err)
		}
	}

	return analytics, nil
}

// GetPartnerDashboard returns dashboard data for partner
func (s *PhotoService) GetPartnerDashboard(partnerID string) (map[string]interface{}, error) {
	dashboard := make(map[string]interface{})

	// Total photos
	var totalPhotos int
	err := s.db.QueryRow("SELECT COUNT(*) FROM photos WHERE partner_id = ?", partnerID).Scan(&totalPhotos)
	if err != nil {
		return nil, err
	}

	// Total views
	var totalViews int
	err = s.db.QueryRow("SELECT COALESCE(SUM(view_count), 0) FROM photos WHERE partner_id = ?", partnerID).Scan(&totalViews)
	if err != nil {
		return nil, err
	}

	// Total shares
	var totalShares int
	err = s.db.QueryRow("SELECT COALESCE(SUM(share_count), 0) FROM photos WHERE partner_id = ?", partnerID).Scan(&totalShares)
	if err != nil {
		return nil, err
	}

	// Recent uploads (last 7 days)
	var recentUploads int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM photos
		WHERE partner_id = ? AND created_at > datetime('now', '-7 days')
	`, partnerID).Scan(&recentUploads)
	if err != nil {
		return nil, err
	}

	dashboard["total_photos"] = totalPhotos
	dashboard["total_views"] = totalViews
	dashboard["total_shares"] = totalShares
	dashboard["recent_uploads"] = recentUploads

	return dashboard, nil
}

// GetPartnerAnalytics returns comprehensive analytics for partner
func (s *PhotoService) GetPartnerAnalytics(partnerID string, from, to *time.Time) (map[string]interface{}, error) {
	analytics := make(map[string]interface{})

	// Build date filter
	dateFilter := ""
	args := []interface{}{partnerID}
	if from != nil && to != nil {
		dateFilter = " AND created_at BETWEEN ? AND ?"
		args = append(args, from, to)
	}

	// Upload trends (daily)
	query := fmt.Sprintf(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM photos WHERE partner_id = ?%s
		GROUP BY DATE(created_at)
		ORDER BY date DESC
		LIMIT 30
	`, dateFilter)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uploadTrends []map[string]interface{}
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			continue
		}
		uploadTrends = append(uploadTrends, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}

	analytics["upload_trends"] = uploadTrends

	return analytics, nil
}