#!/bin/bash

# Create necessary directories
mkdir -p media-vault-backend/internal/auth
mkdir -p media-vault-backend/internal/handlers
mkdir -p media-vault-backend/internal/models
mkdir -p media-vault-backend/internal/services
mkdir -p media-vault-backend/internal/utils

# Extract models
cat > media-vault-backend/internal/models/photo.go << 'EOL'
package models

import "time"

// Photo represents a photo in the vault
type Photo struct {
    ID                string     `json:"id" db:"id"`
    UserID            string     `json:"user_id" db:"user_id"`
    PartnerID         *string    `json:"partner_id,omitempty" db:"partner_id"`
    Filename          string     `json:"filename" db:"filename"`
    OriginalName      string     `json:"original_name" db:"original_name"`
    FilePath          string     `json:"file_path" db:"file_path"`
    ThumbnailPath     *string    `json:"thumbnail_path,omitempty" db:"thumbnail_path"`
    FileSize          int64      `json:"file_size" db:"file_size"`
    MimeType          string     `json:"mime_type" db:"mime_type"`
    Width             *int       `json:"width,omitempty" db:"width"`
    Height            *int       `json:"height,omitempty" db:"height"`
    Hash              string     `json:"hash" db:"hash"`
    Description       *string    `json:"description,omitempty" db:"description"`
    AIDescription     *string    `json:"ai_description,omitempty" db:"ai_description"`
    Tags              *string    `json:"tags,omitempty" db:"tags"`
    AIConfidence      *float64   `json:"ai_confidence,omitempty" db:"ai_confidence"`
    IsNSFW            *bool      `json:"is_nsfw,omitempty" db:"is_nsfw"`
    NSFWConfidence    *float64   `json:"nsfw_confidence,omitempty" db:"nsfw_confidence"`
    ModerationStatus  string     `json:"moderation_status" db:"moderation_status"`
    ExifData          *string    `json:"exif_data,omitempty" db:"exif_data"`
    Location          *string    `json:"location,omitempty" db:"location"`
    CameraMake        *string    `json:"camera_make,omitempty" db:"camera_make"`
    CameraModel       *string    `json:"camera_model,omitempty" db:"camera_model"`
    TakenAt           *time.Time `json:"taken_at,omitempty" db:"taken_at"`
    IsShared          bool       `json:"is_shared" db:"is_shared"`
    SharedWith        []string   `json:"shared_with,omitempty"`
    ShareCount        int        `json:"share_count" db:"share_count"`
    ViewCount         int        `json:"view_count" db:"view_count"`
    CreatedAt         time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
    ProcessedAt       *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// PhotoSharing represents photo sharing permissions
type PhotoSharing struct {
    ID          string     `json:"id" db:"id"`
    PhotoID     string     `json:"photo_id" db:"photo_id"`
    SharedBy    string     `json:"shared_by" db:"shared_by"`
    SharedWith  string     `json:"shared_with" db:"shared_with"`
    Permission  string     `json:"permission" db:"permission"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// PhotoAnalytics represents photo analytics data
type PhotoAnalytics struct {
    PhotoID    string     `json:"photo_id" db:"photo_id"`
    Views      int        `json:"views" db:"views"`
    Downloads  int        `json:"downloads" db:"downloads"`
    Shares     int        `json:"shares" db:"shares"`
    LastViewed *time.Time `json:"last_viewed,omitempty" db:"last_viewed"`
    ViewerIPs  []string   `json:"viewer_ips,omitempty"`
    Countries  []string   `json:"countries,omitempty"`
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
    Operation   string   `json:"operation"`
}

// ShareRequest represents a photo sharing request
type ShareRequest struct {
    PhotoIDs   []string   `json:"photo_ids"`
    ShareWith  []string   `json:"share_with"`
    Permission string     `json:"permission"`
    ExpiresAt  *time.Time `json:"expires_at,omitempty"`
    Message    string     `json:"message,omitempty"`
}
EOL

# Extract partner handler
cat > media-vault-backend/internal/handlers/partner.go << 'EOL'
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"media-vault/internal/services"
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
    // Implementation from main.go
    return c.SendString("Bulk upload endpoint")
}

// GetPartnerPhotos returns photos uploaded by the partner
func (h *PartnerHandler) GetPartnerPhotos(c *fiber.Ctx) error {
    // Implementation from main.go
    return c.SendString("Partner photos endpoint")
}

// BatchUpdateDescriptions updates descriptions for multiple photos
func (h *PartnerHandler) BatchUpdateDescriptions(c *fiber.Ctx) error {
    // Implementation from main.go
    return c.SendString("Batch update descriptions endpoint")
}

// BatchSharePhotos shares multiple photos with users
func (h *PartnerHandler) BatchSharePhotos(c *fiber.Ctx) error {
    // Implementation from main.go
    return c.SendString("Batch share photos endpoint")
}

// GetPhotoAnalytics returns analytics for a specific photo
func (h *PartnerHandler) GetPhotoAnalytics(c *fiber.Ctx) error {
    // Implementation from main.go
    return c.SendString("Photo analytics endpoint")
}

// GetPartnerDashboard returns dashboard data for partner
func (h *PartnerHandler) GetPartnerDashboard(c *fiber.Ctx) error {
    // Implementation from main.go
    return c.SendString("Partner dashboard endpoint")
}

// GetPartnerAnalytics returns comprehensive analytics for partner
func (h *PartnerHandler) GetPartnerAnalytics(c *fiber.Ctx) error {
    // Implementation from main.go
    return c.SendString("Partner analytics endpoint")
}
EOL

# Extract photo service
cat > media-vault-backend/internal/services/photo_service.go << 'EOL'
package services

import (
	"database/sql"
	"mime/multipart"
	"time"
	"media-vault/internal/models"
)

type PhotoService struct {
    db *sql.DB
}

func NewPhotoService(db *sql.DB) *PhotoService {
    return &PhotoService{db: db}
}

// CreatePhotoFromUpload creates a photo record from uploaded file
func (s *PhotoService) CreatePhotoFromUpload(fileHeader *multipart.FileHeader, userID string, partnerID *string) (*models.Photo, error) {
    // Implementation from main.go
    return &models.Photo{}, nil
}

// GetPhoto retrieves a photo by ID
func (s *PhotoService) GetPhoto(photoID string) (*models.Photo, error) {
    // Implementation from main.go
    return &models.Photo{}, nil
}

// UpdatePhoto updates a photo record
func (s *PhotoService) UpdatePhoto(photo *models.Photo) error {
    // Implementation from main.go
    return nil
}

// GetPartnerPhotos retrieves photos for a partner with pagination
func (s *PhotoService) GetPartnerPhotos(partnerID string, page, limit int, sortBy, sortOrder, search string) ([]*models.Photo, int, error) {
    // Implementation from main.go
    return []*models.Photo{}, 0, nil
}

// GenerateAIDescription generates AI description for a photo
func (s *PhotoService) GenerateAIDescription(photoID string) (string, error) {
    // Implementation from main.go
    return "", nil
}

// GetPhotoAnalytics retrieves analytics for a photo
func (s *PhotoService) GetPhotoAnalytics(photoID string) (*models.PhotoAnalytics, error) {
    // Implementation from main.go
    return &models.PhotoAnalytics{}, nil
}

// GetPartnerDashboard returns dashboard data for partner
func (s *PhotoService) GetPartnerDashboard(partnerID string) (map[string]interface{}, error) {
    // Implementation from main.go
    return map[string]interface{}{
        "total_photos": 0,
        "total_views": 0,
        "total_shares": 0,
    }, nil
}

// GetPartnerAnalytics returns comprehensive analytics for partner
func (s *PhotoService) GetPartnerAnalytics(partnerID string, from, to *time.Time) (map[string]interface{}, error) {
    // Implementation from main.go
    return map[string]interface{}{
        "total_photos": 0,
        "total_views": 0,
        "total_shares": 0,
        "views_by_date": []interface{}{},
        "top_photos":   []interface{}{},
    }, nil
}
EOL

echo "Extracted Go files to their respective locations"
