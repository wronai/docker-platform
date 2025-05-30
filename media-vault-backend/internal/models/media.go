package models

import (
	"time"
)

// Media represents a media file in the system
type Media struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	MIMEType    string    `json:"mime_type"`
	Size        int64     `json:"size"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewMedia creates a new Media instance
func NewMedia(userID, title, url, mimeType string, size int64) *Media {
	now := time.Now()
	return &Media{
		ID:        generateMediaID(),
		UserID:    userID,
		Title:     title,
		URL:       url,
		MIMEType:  mimeType,
		Size:      size,
		IsPublic:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the media metadata
func (m *Media) Update(updates map[string]interface{}) {
	if title, ok := updates["title"].(string); ok && title != "" {
		m.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		m.Description = description
	}
	if isPublic, ok := updates["is_public"].(bool); ok {
		m.IsPublic = isPublic
	}
	m.UpdatedAt = time.Now()
}

// generateMediaID generates a new unique media ID
func generateMediaID() string {
	// In a real implementation, you might want to use UUID or another unique ID generator
	return "media_" + time.Now().Format("20060102150405")
}
