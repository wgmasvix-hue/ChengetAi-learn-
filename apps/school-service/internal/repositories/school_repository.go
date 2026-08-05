package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/school-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// SchoolRepository handles school data access
type SchoolRepository struct {
	db *pgxpool.Pool
}

// NewSchoolRepository creates a new school repository
func NewSchoolRepository(db *pgxpool.Pool) *SchoolRepository {
	return &SchoolRepository{db: db}
}

// Create creates a new school
func (r *SchoolRepository) Create(ctx context.Context, school *domain.School) error {
	query := `
		INSERT INTO schools (
			name, description, country, region, city, website, email, phone_number, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		school.Name,
		school.Description,
		school.Country,
		school.Region,
		school.City,
		school.Website,
		school.Email,
		school.PhoneNumber,
		"active",
	).Scan(&school.ID, &school.CreatedAt, &school.UpdatedAt)

	if err != nil {
		return fmt.Errorf("creating school: %w", err)
	}

	return nil
}

// GetByID retrieves a school by ID
func (r *SchoolRepository) GetByID(ctx context.Context, schoolID string) (*domain.School, error) {
	query := `
		SELECT id, name, description, country, region, city, website, email, phone_number,
		       logo_url, status, subscription_tier, subscription_expires_at, created_at, updated_at
		FROM schools
		WHERE id = $1
	`

	school := &domain.School{}
	err := r.db.QueryRow(ctx, query, schoolID).Scan(
		&school.ID,
		&school.Name,
		&school.Description,
		&school.Country,
		&school.Region,
		&school.City,
		&school.Website,
		&school.Email,
		&school.PhoneNumber,
		&school.LogoURL,
		&school.Status,
		&school.SubscriptionTier,
		&school.SubscriptionExpiresAt,
		&school.CreatedAt,
		&school.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewNotFound("school not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("querying school: %w", err)
	}

	return school, nil
}

// Update updates a school
func (r *SchoolRepository) Update(ctx context.Context, school *domain.School) error {
	query := `
		UPDATE schools
		SET name = $1, description = $2, website = $3, email = $4, phone_number = $5, logo_url = $6
		WHERE id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		school.Name,
		school.Description,
		school.Website,
		school.Email,
		school.PhoneNumber,
		school.LogoURL,
		school.ID,
	).Scan(&school.UpdatedAt)

	if err == pgx.ErrNoRows {
		return errors.NewNotFound("school not found", err)
	}
	if err != nil {
		return fmt.Errorf("updating school: %w", err)
	}

	return nil
}

// ListByCountry retrieves schools by country
func (r *SchoolRepository) ListByCountry(ctx context.Context, country string, limit, offset int) ([]domain.School, error) {
	query := `
		SELECT id, name, description, country, region, city, website, email, phone_number,
		       logo_url, status, subscription_tier, subscription_expires_at, created_at, updated_at
		FROM schools
		WHERE country = $1 AND status = 'active'
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, country, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying schools: %w", err)
	}
	defer rows.Close()

	schools := make([]domain.School, 0)
	for rows.Next() {
		school := domain.School{}
		if err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Description,
			&school.Country,
			&school.Region,
			&school.City,
			&school.Website,
			&school.Email,
			&school.PhoneNumber,
			&school.LogoURL,
			&school.Status,
			&school.SubscriptionTier,
			&school.SubscriptionExpiresAt,
			&school.CreatedAt,
			&school.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning school: %w", err)
		}
		schools = append(schools, school)
	}

	return schools, rows.Err()
}
