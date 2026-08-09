package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Repository struct {
	pool  *pgxpool.Pool
	redis *goredis.Client
}

type execQuerier interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type scanQuerier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

func NewRepository(pool *pgxpool.Pool, redisClient *goredis.Client) *Repository {
	return &Repository{pool: pool, redis: redisClient}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.createUser(ctx, r.pool, user)
}

func (r *Repository) CreateUserTx(ctx context.Context, tx pgx.Tx, user *User) error {
	return r.createUser(ctx, tx, user)
}

func (r *Repository) createUser(ctx context.Context, db execQuerier, user *User) error {
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, phone, password_hash, role, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.Email, user.Phone, user.PasswordHash, user.Role, user.FirstName, user.LastName)
	return err
}

func (r *Repository) CreateLearnerProfileTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID, educationLevel, grade, preferredLanguage string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO learner_profiles (user_id, education_level, grade, preferred_language)
		VALUES ($1, $2, $3, $4)
	`, userID, nullableString(educationLevel), nullableString(grade), defaultString(preferredLanguage, "en"))
	return err
}

func (r *Repository) CreateTeacherProfileTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO teacher_profiles (user_id)
		VALUES ($1)
	`, userID)
	return err
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	return r.scanUser(ctx, r.pool, `
		SELECT id, email, phone, password_hash, role, first_name, last_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return r.scanUser(ctx, r.pool, `
		SELECT id, email, phone, password_hash, role, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)
}

func (r *Repository) scanUser(ctx context.Context, db scanQuerier, query string, args ...any) (*User, error) {
	var user User
	if err := db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) StoreSession(ctx context.Context, token string, ttl time.Duration) error {
	if r.redis == nil {
		return errors.New("redis client is nil")
	}
	return r.redis.Set(ctx, sessionKey(token), "1", ttl).Err()
}

func (r *Repository) DeleteSession(ctx context.Context, token string) error {
	if r.redis == nil {
		return errors.New("redis client is nil")
	}
	return r.redis.Del(ctx, sessionKey(token)).Err()
}

func sessionKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "auth:session:" + hex.EncodeToString(sum[:])
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
