package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/recommendations"
	"go.uber.org/zap"
)

// RecommendationService orchestrates recommendation operations
type RecommendationService struct {
	engine *recommendations.Engine
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(
	engine *recommendations.Engine,
	db *pgxpool.Pool,
	logger *zap.SugaredLogger,
) *RecommendationService {
	return &RecommendationService{
		engine: engine,
		db:     db,
		logger: logger,
	}
}

// GetRecommendations retrieves personalized recommendations for a user
func (s *RecommendationService) GetRecommendations(
	ctx context.Context,
	userID string,
	topK int,
	strategy string,
) (*recommendations.RecommendationResult, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if topK == 0 {
		topK = 10
	}

	if strategy == "" {
		strategy = "hybrid"
	}

	s.logger.Debugw("Generating recommendations",
		"user_id", userID,
		"topK", topK,
		"strategy", strategy,
	)

	req := &recommendations.RecommendationRequest{
		UserID:        userID,
		TopK:          topK,
		Strategy:      strategy,
		MinScore:      0.5,
		ExcludeViewed: true,
	}

	result, err := s.engine.GetRecommendations(ctx, req)
	if err != nil {
		s.logger.Errorw("Error generating recommendations", "error", err)
		return nil, err
	}

	s.logger.Infow("Generated recommendations",
		"user_id", userID,
		"count", len(result.Recommendations),
	)

	// Store recommendations in database for tracking
	if err := s.storeRecommendations(ctx, result); err != nil {
		s.logger.Warnw("Error storing recommendations", "error", err)
		// Don't fail if storage fails - recommendations are still valid
	}

	return result, nil
}

// GetMetrics retrieves recommendation effectiveness metrics
func (s *RecommendationService) GetMetrics(
	ctx context.Context,
	userID string,
) map[string]float32 {
	s.logger.Debugw("Retrieving recommendation metrics", "user_id", userID)

	// Generate dummy recommendations for metric calculation
	req := &recommendations.RecommendationRequest{
		UserID:   userID,
		TopK:     10,
		Strategy: "hybrid",
	}

	result, err := s.engine.GetRecommendations(ctx, req)
	if err != nil {
		s.logger.Errorw("Error retrieving metrics", "error", err)
		return make(map[string]float32)
	}

	metrics := s.engine.GetRecommendationMetrics(ctx, userID, result.Recommendations)
	return metrics
}

// SubmitFeedback records user feedback on recommendations
func (s *RecommendationService) SubmitFeedback(
	ctx context.Context,
	userID, resourceID string,
	helpful bool,
	reason string,
) error {
	s.logger.Debugw("Recording recommendation feedback",
		"user_id", userID,
		"resource_id", resourceID,
		"helpful", helpful,
	)

	query := `
		INSERT INTO recommendation_feedback (user_id, resource_id, helpful, reason, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := s.db.Exec(ctx, query, userID, resourceID, helpful, reason)
	if err != nil {
		s.logger.Errorw("Error saving feedback", "error", err)
		return fmt.Errorf("saving feedback: %w", err)
	}

	return nil
}

// storeRecommendations saves generated recommendations to database for tracking
func (s *RecommendationService) storeRecommendations(
	ctx context.Context,
	result *recommendations.RecommendationResult,
) error {
	for i, rec := range result.Recommendations {
		query := `
			INSERT INTO recommendation_history (user_id, resource_id, score, rank, strategy, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`

		_, err := s.db.Exec(ctx, query,
			result.UserID,
			rec.ResourceID,
			rec.Score,
			i+1,
			result.Strategy,
		)

		if err != nil {
			return fmt.Errorf("storing recommendation: %w", err)
		}
	}

	return nil
}
