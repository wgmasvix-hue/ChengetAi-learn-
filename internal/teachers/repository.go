package teachers

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindByUserID(ctx context.Context, id uuid.UUID) (*Profile, error) {
	var profile Profile
	err := r.pool.QueryRow(ctx, `
		SELECT u.id, u.email, u.first_name, u.last_name
		FROM teacher_profiles tp
		JOIN users u ON u.id = tp.user_id
		WHERE tp.user_id = $1
	`, id).Scan(&profile.UserID, &profile.Email, &profile.FirstName, &profile.LastName)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
