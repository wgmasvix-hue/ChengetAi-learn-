package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/auth-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// UserRepository handles user data access
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			email, username, password_hash, full_name, country, language, timezone, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.Country,
		user.Language,
		user.Timezone,
		"active",
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			return errors.NewEmailTaken()
		}
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)" {
			return errors.NewUsernameTaken()
		}
		return fmt.Errorf("creating user: %w", err)
	}

	// Create wallet for new user
	walletQuery := `
		INSERT INTO wallets (user_id, currency)
		VALUES ($1, 'USD')
	`
	if _, err := r.db.Exec(ctx, walletQuery, user.ID); err != nil {
		return fmt.Errorf("creating wallet: %w", err)
	}

	return nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, avatar_url, bio,
		       phone_number, country, language, timezone, status, email_verified,
		       email_verified_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&user.Bio,
		&user.PhoneNumber,
		&user.Country,
		&user.Language,
		&user.Timezone,
		&user.Status,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewNotFound("user not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by email: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, avatar_url, bio,
		       phone_number, country, language, timezone, status, email_verified,
		       email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.AvatarURL,
		&user.Bio,
		&user.PhoneNumber,
		&user.Country,
		&user.Language,
		&user.Timezone,
		&user.Status,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewNotFound("user not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("querying user by id: %w", err)
	}

	return user, nil
}
