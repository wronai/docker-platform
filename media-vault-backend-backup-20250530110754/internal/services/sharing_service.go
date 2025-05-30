package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Share represents a shared photo with permissions
type Share struct {
	ID          string    `json:"id"`
	PhotoID     string    `json:"photo_id"`
	SharedBy    string    `json:"shared_by"`
	SharedWith  string    `json:"shared_with"`
	Permission  string    `json:"permission"` // "view" or "download"
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// SharingService handles photo sharing operations
type SharingService struct {
	db *sql.DB
}

// NewSharingService creates a new SharingService
func NewSharingService(db *sql.DB) *SharingService {
	return &SharingService{
		db: db,
	}
}

// SharePhoto shares a photo with another user
func (s *SharingService) SharePhoto(ctx context.Context, share *Share) error {
	// Validate input
	if share.PhotoID == "" || share.SharedBy == "" || share.SharedWith == "" {
		return errors.New("photo ID, shared_by, and shared_with are required")
	}

	// Set default permission if not provided
	if share.Permission == "" {
		share.Permission = "view"
	}

	// Generate a new share ID if not provided
	if share.ID == "" {
		share.ID = uuid.New().String()
	}

	// Set created at time
	share.CreatedAt = time.Now()

	// Insert the share record
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO photo_shares (id, photo_id, shared_by, shared_with, permission, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		share.ID,
		share.PhotoID,
		share.SharedBy,
		share.SharedWith,
		share.Permission,
		share.ExpiresAt,
		share.CreatedAt,
	)

	return err
}

// GetShare retrieves a share by ID
func (s *SharingService) GetShare(ctx context.Context, shareID string) (*Share, error) {
	var share Share
	err := s.db.QueryRowContext(ctx, `
		SELECT id, photo_id, shared_by, shared_with, permission, expires_at, created_at
		FROM photo_shares
		WHERE id = $1
	`, shareID).Scan(
		&share.ID,
		&share.PhotoID,
		&share.SharedBy,
		&share.SharedWith,
		&share.Permission,
		&share.ExpiresAt,
		&share.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("share not found")
		}
		return nil, err
	}

	return &share, nil
}

// ListSharesForPhoto lists all shares for a specific photo
func (s *SharingService) ListSharesForPhoto(ctx context.Context, photoID string) ([]*Share, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, photo_id, shared_by, shared_with, permission, expires_at, created_at
		FROM photo_shares
		WHERE photo_id = $1
	`, photoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []*Share
	for rows.Next() {
		var share Share
		if err := rows.Scan(
			&share.ID,
			&share.PhotoID,
			&share.SharedBy,
			&share.SharedWith,
			&share.Permission,
			&share.ExpiresAt,
			&share.CreatedAt,
		); err != nil {
			return nil, err
		}
		shares = append(shares, &share)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shares, nil
}

// RevokeShare removes a share
func (s *SharingService) RevokeShare(ctx context.Context, shareID string) error {
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM photo_shares
		WHERE id = $1
	`, shareID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("share not found")
	}

	return nil
}

// HasPermission checks if a user has permission to access a photo
func (s *SharingService) HasPermission(ctx context.Context, photoID, userID, requiredPermission string) (bool, error) {
	// Check if the user is the owner of the photo
	var ownerID string
	err := s.db.QueryRowContext(ctx, `
		SELECT user_id FROM photos WHERE id = $1
	`, photoID).Scan(&ownerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("photo not found")
		}
		return false, err
	}

	// If the user is the owner, they have all permissions
	if ownerID == userID {
		return true, nil
	}

	// Check for an active share
	var count int
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM photo_shares
		WHERE photo_id = $1
		  AND shared_with = $2
		  AND (expires_at IS NULL OR expires_at > NOW())
		  AND permission = $3
	`, photoID, userID, requiredPermission).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}