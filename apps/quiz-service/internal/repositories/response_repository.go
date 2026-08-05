package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/quiz"
)

// ResponseRepository handles quiz response data access
type ResponseRepository struct {
	db *pgxpool.Pool
}

// NewResponseRepository creates a new response repository
func NewResponseRepository(db *pgxpool.Pool) *ResponseRepository {
	return &ResponseRepository{db: db}
}

// Create saves a quiz response
func (r *ResponseRepository) Create(ctx context.Context, response *quiz.QuizResponse) error {
	responsesJSON, err := json.Marshal(response.Responses)
	if err != nil {
		return fmt.Errorf("marshaling responses: %w", err)
	}

	query := `
		INSERT INTO quiz_responses (id, quiz_id, user_id, responses, score, percentage, time_spent, submitted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.Exec(ctx, query,
		response.ID,
		response.QuizID,
		response.UserID,
		responsesJSON,
		response.Score,
		response.Percentage,
		response.TimeSpent,
		response.SubmittedAt,
	)

	if err != nil {
		return fmt.Errorf("creating response: %w", err)
	}

	return nil
}

// GetByID retrieves a quiz response
func (r *ResponseRepository) GetByID(ctx context.Context, responseID string) (*quiz.QuizResponse, error) {
	query := `
		SELECT id, quiz_id, user_id, responses, score, percentage, time_spent, submitted_at
		FROM quiz_responses
		WHERE id = $1
	`

	var resp quiz.QuizResponse
	var responsesJSON []byte

	err := r.db.QueryRow(ctx, query, responseID).Scan(
		&resp.ID,
		&resp.QuizID,
		&resp.UserID,
		&responsesJSON,
		&resp.Score,
		&resp.Percentage,
		&resp.TimeSpent,
		&resp.SubmittedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying response: %w", err)
	}

	if err := json.Unmarshal(responsesJSON, &resp.Responses); err != nil {
		return nil, fmt.Errorf("unmarshaling responses: %w", err)
	}

	return &resp, nil
}

// GetByUserAndQuiz retrieves responses for a specific user and quiz
func (r *ResponseRepository) GetByUserAndQuiz(ctx context.Context, userID, quizID string) ([]quiz.QuizResponse, error) {
	query := `
		SELECT id, quiz_id, user_id, responses, score, percentage, time_spent, submitted_at
		FROM quiz_responses
		WHERE user_id = $1 AND quiz_id = $2
		ORDER BY submitted_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, quizID)
	if err != nil {
		return nil, fmt.Errorf("querying responses: %w", err)
	}
	defer rows.Close()

	responses := make([]quiz.QuizResponse, 0)
	for rows.Next() {
		var resp quiz.QuizResponse
		var responsesJSON []byte

		if err := rows.Scan(
			&resp.ID,
			&resp.QuizID,
			&resp.UserID,
			&responsesJSON,
			&resp.Score,
			&resp.Percentage,
			&resp.TimeSpent,
			&resp.SubmittedAt,
		); err != nil {
			continue
		}

		if err := json.Unmarshal(responsesJSON, &resp.Responses); err != nil {
			continue
		}

		responses = append(responses, resp)
	}

	return responses, rows.Err()
}

// GetByUser retrieves all responses for a user
func (r *ResponseRepository) GetByUser(ctx context.Context, userID string, limit, offset int) ([]quiz.QuizResponse, error) {
	query := `
		SELECT id, quiz_id, user_id, responses, score, percentage, time_spent, submitted_at
		FROM quiz_responses
		WHERE user_id = $1
		ORDER BY submitted_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("querying responses: %w", err)
	}
	defer rows.Close()

	responses := make([]quiz.QuizResponse, 0)
	for rows.Next() {
		var resp quiz.QuizResponse
		var responsesJSON []byte

		if err := rows.Scan(
			&resp.ID,
			&resp.QuizID,
			&resp.UserID,
			&responsesJSON,
			&resp.Score,
			&resp.Percentage,
			&resp.TimeSpent,
			&resp.SubmittedAt,
		); err != nil {
			continue
		}

		if err := json.Unmarshal(responsesJSON, &resp.Responses); err != nil {
			continue
		}

		responses = append(responses, resp)
	}

	return responses, rows.Err()
}

// GetStats returns statistics for a user on a quiz
type QuizStats struct {
	TotalAttempts  int
	AverageScore   float32
	HighestScore   int
	LastAttempted  int64
	PassingRate    float32
}

// GetUserStats retrieves statistics for a user's quiz attempts
func (r *ResponseRepository) GetUserStats(ctx context.Context, userID, quizID string) (*QuizStats, error) {
	query := `
		SELECT COUNT(*) as total_attempts,
		       AVG(percentage) as avg_percentage,
		       MAX(score) as max_score,
		       MAX(submitted_at) as last_attempted,
		       SUM(CASE WHEN percentage >= 70 THEN 1 ELSE 0 END)::float / COUNT(*) as pass_rate
		FROM quiz_responses
		WHERE user_id = $1 AND quiz_id = $2
	`

	var stats QuizStats
	var avgPercentage, passRate *float32
	var maxScore *int
	var lastAttempted *int64

	err := r.db.QueryRow(ctx, query, userID, quizID).Scan(
		&stats.TotalAttempts,
		&avgPercentage,
		&maxScore,
		&lastAttempted,
		&passRate,
	)

	if err != nil {
		return nil, fmt.Errorf("querying stats: %w", err)
	}

	if avgPercentage != nil {
		stats.AverageScore = *avgPercentage
	}
	if maxScore != nil {
		stats.HighestScore = *maxScore
	}
	if lastAttempted != nil {
		stats.LastAttempted = *lastAttempted
	}
	if passRate != nil {
		stats.PassingRate = *passRate
	}

	return &stats, nil
}
