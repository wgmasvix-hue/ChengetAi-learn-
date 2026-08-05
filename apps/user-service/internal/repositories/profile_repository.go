package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/user-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// ProfileRepository handles profile data access
type ProfileRepository struct {
	db *pgxpool.Pool
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetByUserID retrieves a profile by user ID
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	query := `
		SELECT id, email, username, full_name, avatar_url, bio, phone_number,
		       country, language, timezone, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	profile := &domain.Profile{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&userID, // Placeholder for unused id
		&userID, // Placeholder for unused email
		&userID, // Placeholder for unused username
		&profile.FullName,
		&profile.AvatarURL,
		&profile.Bio,
		&profile.PhoneNumber,
		&profile.Country,
		&profile.Language,
		&profile.Timezone,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewNotFound("user not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("querying profile: %w", err)
	}

	profile.UserID = userID
	return profile, nil
}

// Update updates a profile
func (r *ProfileRepository) Update(ctx context.Context, profile *domain.Profile) error {
	query := `
		UPDATE users
		SET full_name = $1, avatar_url = $2, bio = $3, phone_number = $4, timezone = $5
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		profile.FullName,
		profile.AvatarURL,
		profile.Bio,
		profile.PhoneNumber,
		profile.Timezone,
		profile.UserID,
	).Scan(&profile.UpdatedAt)

	if err == pgx.ErrNoRows {
		return errors.NewNotFound("user not found", err)
	}
	if err != nil {
		return fmt.Errorf("updating profile: %w", err)
	}

	return nil
}

// Delete deletes a profile
func (r *ProfileRepository) Delete(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP, status = 'deleted'
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("deleting profile: %w", err)
	}

	return nil
}
