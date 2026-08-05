package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// Resource represents an educational resource
type Resource struct {
	ID               string
	DSpaceID         string
	Title            string
	Description      *string
	CreatorID        string
	ResourceType     *string
	Curriculum       *string
	Subject          *string
	Country          *string
	Language         string
	GradeLevel       *string
	FileURL          *string
	MimeType         *string
	ViewsCount       int
	DownloadsCount   int
	BookmarksCount   int
	Status           string
}

// ResourceRepository handles resource data access
type ResourceRepository struct {
	db *pgxpool.Pool
}

// NewResourceRepository creates a new resource repository
func NewResourceRepository(db *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// Create creates a new resource
func (r *ResourceRepository) Create(ctx context.Context, resource *Resource) error {
	query := `
		INSERT INTO resources (
			dspace_id, title, description, creator_id, resource_type,
			curriculum, subject, country, language, grade_level,
			file_url, mime_type, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		resource.DSpaceID,
		resource.Title,
		resource.Description,
		resource.CreatorID,
		resource.ResourceType,
		resource.Curriculum,
		resource.Subject,
		resource.Country,
		resource.Language,
		resource.GradeLevel,
		resource.FileURL,
		resource.MimeType,
		"published",
	).Scan(&resource.ID)

	if err != nil {
		return fmt.Errorf("creating resource: %w", err)
	}

	return nil
}

// GetByDSpaceID retrieves a resource by DSpace ID
func (r *ResourceRepository) GetByDSpaceID(ctx context.Context, dspaceID string) (*Resource, error) {
	query := `
		SELECT id, dspace_id, title, description, creator_id, resource_type,
		       curriculum, subject, country, language, grade_level, file_url,
		       mime_type, views_count, downloads_count, bookmarks_count, status
		FROM resources
		WHERE dspace_id = $1
	`

	resource := &Resource{}
	err := r.db.QueryRow(ctx, query, dspaceID).Scan(
		&resource.ID,
		&resource.DSpaceID,
		&resource.Title,
		&resource.Description,
		&resource.CreatorID,
		&resource.ResourceType,
		&resource.Curriculum,
		&resource.Subject,
		&resource.Country,
		&resource.Language,
		&resource.GradeLevel,
		&resource.FileURL,
		&resource.MimeType,
		&resource.ViewsCount,
		&resource.DownloadsCount,
		&resource.BookmarksCount,
		&resource.Status,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.NewNotFound("resource not found", err)
	}
	if err != nil {
		return nil, fmt.Errorf("querying resource: %w", err)
	}

	return resource, nil
}

// IncrementViews increments the view count for a resource
func (r *ResourceRepository) IncrementViews(ctx context.Context, resourceID string) error {
	query := `UPDATE resources SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, resourceID)
	if err != nil {
		return fmt.Errorf("incrementing views: %w", err)
	}
	return nil
}

// IncrementDownloads increments the download count for a resource
func (r *ResourceRepository) IncrementDownloads(ctx context.Context, resourceID string) error {
	query := `UPDATE resources SET downloads_count = downloads_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, resourceID)
	if err != nil {
		return fmt.Errorf("incrementing downloads: %w", err)
	}
	return nil
}

// SearchByTitle searches resources by title
func (r *ResourceRepository) SearchByTitle(ctx context.Context, query string, limit, offset int) ([]Resource, error) {
	sql := `
		SELECT id, dspace_id, title, description, creator_id, resource_type,
		       curriculum, subject, country, language, grade_level, file_url,
		       mime_type, views_count, downloads_count, bookmarks_count, status
		FROM resources
		WHERE title ILIKE $1 AND status = 'published'
		ORDER BY views_count DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, sql, "%"+query+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("searching resources: %w", err)
	}
	defer rows.Close()

	resources := make([]Resource, 0)
	for rows.Next() {
		resource := Resource{}
		if err := rows.Scan(
			&resource.ID,
			&resource.DSpaceID,
			&resource.Title,
			&resource.Description,
			&resource.CreatorID,
			&resource.ResourceType,
			&resource.Curriculum,
			&resource.Subject,
			&resource.Country,
			&resource.Language,
			&resource.GradeLevel,
			&resource.FileURL,
			&resource.MimeType,
			&resource.ViewsCount,
			&resource.DownloadsCount,
			&resource.BookmarksCount,
			&resource.Status,
		); err != nil {
			return nil, fmt.Errorf("scanning resource: %w", err)
		}
		resources = append(resources, resource)
	}

	return resources, rows.Err()
}
