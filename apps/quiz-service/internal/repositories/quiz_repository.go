package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/quiz"
)

// QuizRepository handles quiz data access
type QuizRepository struct {
	db *pgxpool.Pool
}

// NewQuizRepository creates a new quiz repository
func NewQuizRepository(db *pgxpool.Pool) *QuizRepository {
	return &QuizRepository{db: db}
}

// Create saves a new quiz
func (r *QuizRepository) Create(ctx context.Context, q *quiz.Quiz) error {
	questionsJSON, err := json.Marshal(q.Questions)
	if err != nil {
		return fmt.Errorf("marshaling questions: %w", err)
	}

	query := `
		INSERT INTO quizzes (id, resource_id, title, description, difficulty_level,
		                      estimated_time, passing_score, questions, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(ctx, query,
		q.ID,
		q.ResourceID,
		q.Title,
		q.Description,
		q.DifficultyLevel,
		q.EstimatedTime,
		q.PassingScore,
		questionsJSON,
		q.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("creating quiz: %w", err)
	}

	return nil
}

// GetByID retrieves a quiz by ID
func (r *QuizRepository) GetByID(ctx context.Context, quizID string) (*quiz.Quiz, error) {
	query := `
		SELECT id, resource_id, title, description, difficulty_level,
		       estimated_time, passing_score, questions, created_at
		FROM quizzes
		WHERE id = $1
	`

	var q quiz.Quiz
	var questionsJSON []byte

	err := r.db.QueryRow(ctx, query, quizID).Scan(
		&q.ID,
		&q.ResourceID,
		&q.Title,
		&q.Description,
		&q.DifficultyLevel,
		&q.EstimatedTime,
		&q.PassingScore,
		&questionsJSON,
		&q.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying quiz: %w", err)
	}

	if err := json.Unmarshal(questionsJSON, &q.Questions); err != nil {
		return nil, fmt.Errorf("unmarshaling questions: %w", err)
	}

	return &q, nil
}

// GetByResourceID retrieves quizzes for a resource
func (r *QuizRepository) GetByResourceID(ctx context.Context, resourceID string) ([]quiz.Quiz, error) {
	query := `
		SELECT id, resource_id, title, description, difficulty_level,
		       estimated_time, passing_score, questions, created_at
		FROM quizzes
		WHERE resource_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, resourceID)
	if err != nil {
		return nil, fmt.Errorf("querying quizzes: %w", err)
	}
	defer rows.Close()

	quizzes := make([]quiz.Quiz, 0)
	for rows.Next() {
		var q quiz.Quiz
		var questionsJSON []byte

		if err := rows.Scan(
			&q.ID,
			&q.ResourceID,
			&q.Title,
			&q.Description,
			&q.DifficultyLevel,
			&q.EstimatedTime,
			&q.PassingScore,
			&questionsJSON,
			&q.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning quiz: %w", err)
		}

		if err := json.Unmarshal(questionsJSON, &q.Questions); err != nil {
			continue // Skip quizzes with corrupted questions
		}

		quizzes = append(quizzes, q)
	}

	return quizzes, rows.Err()
}

// Delete removes a quiz
func (r *QuizRepository) Delete(ctx context.Context, quizID string) error {
	query := `DELETE FROM quizzes WHERE id = $1`
	_, err := r.db.Exec(ctx, query, quizID)
	if err != nil {
		return fmt.Errorf("deleting quiz: %w", err)
	}
	return nil
}
