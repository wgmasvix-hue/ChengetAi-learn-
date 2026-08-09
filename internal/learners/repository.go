package learners

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	UserID            string `json:"user_id"`
	Email             string `json:"email"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	EducationLevel    string `json:"education_level,omitempty"`
	Grade             string `json:"grade,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
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
		SELECT u.id, u.email, u.first_name, u.last_name, lp.education_level, lp.grade, lp.preferred_language
		FROM learner_profiles lp
		JOIN users u ON u.id = lp.user_id
		WHERE lp.user_id = $1
	`, id).Scan(
		&profile.UserID,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&profile.EducationLevel,
		&profile.Grade,
		&profile.PreferredLanguage,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
