package models

import (
	"time"
)

// Description represents a description for a media item
type Description struct {
	ID          string    `json:"id"`
	MediaID     string    `json:"media_id"`
	Content     string    `json:"content"`
	GeneratedBy string    `json:"generated_by"` // "user" or "ai"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewDescription creates a new Description instance
func NewDescription(mediaID, content, generatedBy string) *Description {
	now := time.Now()
	return &Description{
		ID:          generateID(),
		MediaID:     mediaID,
		Content:     content,
		GeneratedBy: generatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the description content
func (d *Description) Update(content string) {
	d.Content = content
	d.UpdatedAt = time.Now()
}

// generateID generates a new unique ID
func generateID() string {
	// In a real implementation, you might want to use UUID or another unique ID generator
	return "desc_" + time.Now().Format("20060102150405")
}