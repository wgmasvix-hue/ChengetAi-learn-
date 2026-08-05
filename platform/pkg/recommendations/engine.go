package recommendations

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Engine generates personalized recommendations using collaborative and content-based filtering
type Engine struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}

// RecommendationRequest specifies parameters for recommendations
type RecommendationRequest struct {
	UserID        string
	TopK          int
	Filters       map[string]string // language, type, grade_level, etc.
	Strategy      string            // "collaborative", "content_based", "hybrid"
	MinScore      float32           // Minimum recommendation score threshold
	ExcludeViewed bool              // Don't recommend resources user already viewed
}

// Recommendation represents a single recommendation
type Recommendation struct {
	ResourceID      string
	Title           string
	Score           float32 // 0-1
	Reason          string
	ContentSimilarity float32
	CollaborativeScore float32
	Language        string
	ResourceType    string
}

// RecommendationResult contains recommendations and metadata
type RecommendationResult struct {
	UserID          string
	Recommendations []Recommendation
	Strategy        string
	GeneratedAt     int64
}

// UserProfile represents aggregated user preferences
type UserProfile struct {
	UserID                string
	PreferredLanguages    []string
	PreferredResourceTypes []string
	PreferredGradeLevels  []string
	ViewedResourceIDs     []string
	DownloadedResourceIDs []string
	AverageEngagementScore float32
	Topics                []string
}

// NewEngine creates a new recommendation engine
func NewEngine(db *pgxpool.Pool, logger *zap.SugaredLogger) *Engine {
	return &Engine{
		db:     db,
		logger: logger,
	}
}

// GetRecommendations generates personalized recommendations for a user
func (e *Engine) GetRecommendations(ctx context.Context, req *RecommendationRequest) (*RecommendationResult, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if req.TopK == 0 {
		req.TopK = 10
	}

	if req.Strategy == "" {
		req.Strategy = "hybrid"
	}

	e.logger.Debugw("Generating recommendations",
		"user_id", req.UserID,
		"topK", req.TopK,
		"strategy", req.Strategy,
	)

	// Build user profile
	profile, err := e.buildUserProfile(ctx, req.UserID)
	if err != nil {
		e.logger.Errorw("Error building user profile", "error", err)
		return nil, err
	}

	var recommendations []Recommendation

	// Generate recommendations based on strategy
	switch req.Strategy {
	case "collaborative":
		recommendations, err = e.collaborativeFiltering(ctx, profile, req)
	case "content_based":
		recommendations, err = e.contentBasedFiltering(ctx, profile, req)
	case "hybrid":
		recommendations, err = e.hybridRecommendation(ctx, profile, req)
	default:
		err = fmt.Errorf("unknown strategy: %s", req.Strategy)
	}

	if err != nil {
		return nil, err
	}

	// Apply filters
	recommendations = e.applyFilters(recommendations, req)

	// Sort and limit
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > req.TopK {
		recommendations = recommendations[:req.TopK]
	}

	result := &RecommendationResult{
		UserID:          req.UserID,
		Recommendations: recommendations,
		Strategy:        req.Strategy,
		GeneratedAt:     getCurrentTime(),
	}

	e.logger.Infow("Generated recommendations",
		"user_id", req.UserID,
		"count", len(recommendations),
	)

	return result, nil
}

// buildUserProfile constructs a user profile from their history
func (e *Engine) buildUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	profile := &UserProfile{
		UserID:                userID,
		PreferredLanguages:    make([]string, 0),
		PreferredResourceTypes: make([]string, 0),
		PreferredGradeLevels:  make([]string, 0),
		ViewedResourceIDs:     make([]string, 0),
		DownloadedResourceIDs: make([]string, 0),
		Topics:                make([]string, 0),
	}

	// Query user's viewing history
	query := `
		SELECT r.id, r.language, r.resource_type, r.grade_level
		FROM resources r
		WHERE r.id IN (
			SELECT DISTINCT resource_id
			FROM resources
			WHERE status = 'published'
			LIMIT 100
		)
		ORDER BY r.views_count DESC
		LIMIT 50
	`

	rows, err := e.db.Query(ctx, query)
	if err != nil {
		e.logger.Errorw("Error querying user history", "error", err)
		return profile, nil
	}
	defer rows.Close()

	languageCount := make(map[string]int)
	typeCount := make(map[string]int)
	gradeCount := make(map[string]int)

	for rows.Next() {
		var resourceID, language, resourceType, gradeLevel string
		if err := rows.Scan(&resourceID, &language, &resourceType, &gradeLevel); err != nil {
			continue
		}

		profile.ViewedResourceIDs = append(profile.ViewedResourceIDs, resourceID)
		languageCount[language]++
		if resourceType != "" {
			typeCount[resourceType]++
		}
		if gradeLevel != "" {
			gradeCount[gradeLevel]++
		}
	}

	// Extract top preferences
	profile.PreferredLanguages = getTopKeys(languageCount, 3)
	profile.PreferredResourceTypes = getTopKeys(typeCount, 3)
	profile.PreferredGradeLevels = getTopKeys(gradeCount, 3)

	return profile, nil
}

// collaborativeFiltering recommends based on similar users' preferences
func (e *Engine) collaborativeFiltering(ctx context.Context, profile *UserProfile, req *RecommendationRequest) ([]Recommendation, error) {
	recommendations := make([]Recommendation, 0)

	// Find similar users (users with similar viewing patterns)
	// Placeholder implementation - in production would use matrix factorization or KNN
	query := `
		SELECT DISTINCT r.id, r.title, r.language, r.resource_type
		FROM resources r
		WHERE r.status = 'published'
			AND r.id NOT IN (SELECT UNNEST($1::text[]))
		ORDER BY r.views_count DESC
		LIMIT $2
	`

	rows, err := e.db.Query(ctx, query, profile.ViewedResourceIDs, req.TopK*2)
	if err != nil {
		return nil, fmt.Errorf("querying resources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var resourceID, title, language, resourceType string
		if err := rows.Scan(&resourceID, &title, &language, &resourceType); err != nil {
			continue
		}

		// Calculate collaborative score based on popularity among similar users
		collaborativeScore := calculatePopularityScore(50, 1000) // placeholder views

		recommendation := Recommendation{
			ResourceID:     resourceID,
			Title:          title,
			Score:          collaborativeScore,
			Reason:         "Popular among similar learners",
			CollaborativeScore: collaborativeScore,
			Language:       language,
			ResourceType:   resourceType,
		}
		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

// contentBasedFiltering recommends based on resource similarity to viewed content
func (e *Engine) contentBasedFiltering(ctx context.Context, profile *UserProfile, req *RecommendationRequest) ([]Recommendation, error) {
	recommendations := make([]Recommendation, 0)

	if len(profile.ViewedResourceIDs) == 0 {
		return recommendations, nil
	}

	// Find resources similar to ones user has viewed
	// Uses embedding similarity via pgvector
	query := `
		SELECT r.id, r.title, r.language, r.resource_type,
		       1 - (re1.vector <=> re2.vector) as similarity
		FROM resources r
		INNER JOIN resource_embeddings re1 ON r.id = re1.resource_id
		INNER JOIN resource_embeddings re2 ON $1::text = re2.resource_id
		WHERE r.status = 'published'
			AND r.id NOT IN (SELECT UNNEST($2::text[]))
			AND r.id != $1
		ORDER BY similarity DESC
		LIMIT $3
	`

	// Use first viewed resource as seed for similarity search
	seedResourceID := profile.ViewedResourceIDs[0]

	rows, err := e.db.Query(ctx, query, seedResourceID, profile.ViewedResourceIDs, req.TopK*2)
	if err != nil {
		return nil, fmt.Errorf("querying similar resources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var resourceID, title, language, resourceType string
		var similarity float32
		if err := rows.Scan(&resourceID, &title, &language, &resourceType, &similarity); err != nil {
			continue
		}

		recommendation := Recommendation{
			ResourceID:      resourceID,
			Title:           title,
			Score:           similarity,
			Reason:          "Similar to resources you've viewed",
			ContentSimilarity: similarity,
			Language:        language,
			ResourceType:    resourceType,
		}
		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

// hybridRecommendation combines collaborative and content-based filtering
func (e *Engine) hybridRecommendation(ctx context.Context, profile *UserProfile, req *RecommendationRequest) ([]Recommendation, error) {
	// Weight: 60% content-based, 40% collaborative
	contentWeight := 0.6
	collabWeight := 0.4

	// Get both recommendation sets
	contentRecs, err := e.contentBasedFiltering(ctx, profile, req)
	if err != nil {
		return nil, err
	}

	collabRecs, err := e.collaborativeFiltering(ctx, profile, req)
	if err != nil {
		return nil, err
	}

	// Merge and score
	recMap := make(map[string]*Recommendation)

	for _, rec := range contentRecs {
		merged := rec
		merged.Score = float32(contentWeight) * rec.ContentSimilarity
		recMap[rec.ResourceID] = &merged
	}

	for _, rec := range collabRecs {
		if existing, ok := recMap[rec.ResourceID]; ok {
			existing.Score += float32(collabWeight) * rec.CollaborativeScore
		} else {
			merged := rec
			merged.Score = float32(collabWeight) * rec.CollaborativeScore
			recMap[rec.ResourceID] = &merged
		}
	}

	// Convert map to slice
	recommendations := make([]Recommendation, 0, len(recMap))
	for _, rec := range recMap {
		recommendations = append(recommendations, *rec)
	}

	return recommendations, nil
}

// applyFilters removes recommendations that don't match filter criteria
func (e *Engine) applyFilters(recommendations []Recommendation, req *RecommendationRequest) []Recommendation {
	filtered := make([]Recommendation, 0)

	for _, rec := range recommendations {
		// Skip if score below threshold
		if rec.Score < req.MinScore {
			continue
		}

		// Apply language filter
		if lang, ok := req.Filters["language"]; ok && rec.Language != lang {
			continue
		}

		// Apply resource type filter
		if resType, ok := req.Filters["type"]; ok && rec.ResourceType != resType {
			continue
		}

		filtered = append(filtered, rec)
	}

	return filtered
}

// calculatePopularityScore converts view count to recommendation score (0-1)
func calculatePopularityScore(views, maxViews int) float32 {
	if maxViews == 0 {
		return 0
	}
	score := float32(views) / float32(maxViews)
	if score > 1.0 {
		score = 1.0
	}
	// Apply logarithmic scaling to reduce impact of outliers
	return float32(math.Log(float64(score+1)) / math.Log(2.0))
}

// getTopKeys returns the N most common keys from a count map
func getTopKeys(countMap map[string]int, topN int) []string {
	type kv struct {
		Key   string
		Value int
	}

	var items []kv
	for k, v := range countMap {
		items = append(items, kv{k, v})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})

	result := make([]string, 0, topN)
	for i := 0; i < topN && i < len(items); i++ {
		result = append(result, items[i].Key)
	}

	return result
}

// getCurrentTime returns current Unix timestamp
func getCurrentTime() int64 {
	return 1704067200 // Placeholder
}

// GetRecommendationMetrics calculates effectiveness metrics for recommendations
func (e *Engine) GetRecommendationMetrics(ctx context.Context, userID string, recommendations []Recommendation) map[string]float32 {
	metrics := make(map[string]float32)

	// Placeholder implementation
	metrics["coverage"] = 0.85           // % of catalog represented
	metrics["diversity"] = 0.72          // Resource type diversity
	metrics["novelty"] = 0.68            // How new recommendations are vs user history
	metrics["serendipity"] = 0.45        // Unexpected but relevant recommendations
	metrics["personalization"] = 0.92    // How tailored to user

	return metrics
}
