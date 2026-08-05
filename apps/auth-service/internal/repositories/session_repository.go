package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// Session represents a user session
type Session struct {
	ID                string
	UserID            string
	TokenHash         string
	RefreshTokenHash  string
	IPAddress         *string
	UserAgent         *string
	ExpiresAt         time.Time
	RefreshExpiresAt  *time.Time
	LastActivityAt    time.Time
	CreatedAt         time.Time
	RevokedAt         *time.Time
}

// SessionRepository handles session data access
type SessionRepository struct {
	db *pgxpool.Pool
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (
			user_id, token_hash, refresh_token_hash, ip_address, user_agent,
			expires_at, refresh_expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		session.UserID,
		session.TokenHash,
		session.RefreshTokenHash,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
		session.RefreshExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)

	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	return nil
}

// GetByTokenHash retrieves a session by token hash
func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
	query := `
		SELECT id, user_id, token_hash, refresh_token_hash, ip_address, user_agent,
		       expires_at, refresh_expires_at, last_activity_at, created_at, revoked_at
		FROM sessions
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	session := &Session{}
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.RefreshTokenHash,
		&session.IPAddress,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.RefreshExpiresAt,
		&session.LastActivityAt,
		&session.CreatedAt,
		&session.RevokedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewUnauthorized("session not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("querying session: %w", err)
	}

	// Check if session expired
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.NewSessionExpired()
	}

	return session, nil
}

// UpdateLastActivity updates the last activity timestamp
func (r *SessionRepository) UpdateLastActivity(ctx context.Context, sessionID string) error {
	query := `
		UPDATE sessions
		SET last_activity_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("updating session activity: %w", err)
	}

	return nil
}

// Revoke revokes a session
func (r *SessionRepository) Revoke(ctx context.Context, sessionID string) error {
	query := `
		UPDATE sessions
		SET revoked_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("revoking session: %w", err)
	}

	return nil
}
