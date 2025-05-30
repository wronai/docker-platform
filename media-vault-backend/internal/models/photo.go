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
