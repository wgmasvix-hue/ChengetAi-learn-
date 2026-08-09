package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, phone, role, first_name, last_name
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Email, &user.Phone, &user.Role, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
