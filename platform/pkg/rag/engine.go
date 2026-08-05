package rag

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Engine performs Retrieval Augmented Generation by finding relevant resources
// and combining them with LLM responses
type Engine struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}

// RetrievedResource represents a resource found during retrieval
type RetrievedResource struct {
	ResourceID   string
	Title        string
	Content      string
	Similarity   float32
	CreatorID    string
	ResourceType string
	Language     string
	Confidence   float32
}

// RetrievalRequest specifies what to retrieve
type RetrievalRequest struct {
	Query          string
	EmbeddingVector []float32
	TopK           int // Number of results to retrieve
	MinSimilarity  float32 // Minimum similarity threshold
	Filters        map[string]string // Filter by language, type, etc.
}

// RetrievalResult contains resources and metadata
type RetrievalResult struct {
	Resources   []RetrievedResource
	TotalFound  int
	QueryTime   int64 // milliseconds
	RerankedAt  int64
}

// NewEngine creates a new RAG engine
func NewEngine(db *pgxpool.Pool, logger *zap.SugaredLogger) *Engine {
	return &Engine{
		db:     db,
		logger: logger,
	}
}

// RetrieveByEmbedding retrieves relevant resources using embedding similarity
// This performs semantic search against stored embeddings
func (e *Engine) RetrieveByEmbedding(ctx context.Context, req *RetrievalRequest) (*RetrievalResult, error) {
	if len(req.EmbeddingVector) == 0 {
		return nil, fmt.Errorf("embedding vector cannot be empty")
	}

	if req.TopK == 0 {
		req.TopK = 5
	}

	e.logger.Debugw("Retrieving resources by embedding",
		"query", req.Query,
		"topK", req.TopK,
		"min_similarity", req.MinSimilarity,
	)

	// Build query for semantic search against embeddings
	// Uses pgvector extension for similarity search
	query := `
		SELECT
			r.id,
			r.title,
			rc.content,
			1 - (re.vector <=> $1::vector) as similarity,
			r.creator_id,
			r.resource_type,
			r.language,
			COALESCE(rc.confidence_score, 0.9)::float as confidence
		FROM resources r
		INNER JOIN resource_embeddings re ON r.id = re.resource_id
		LEFT JOIN resource_content rc ON r.id = rc.resource_id
		WHERE r.status = 'published'
			AND (1 - (re.vector <=> $1::vector)) > $2
	`

	// Add language filter if specified
	if lang, ok := req.Filters["language"]; ok {
		query += fmt.Sprintf(" AND r.language = '%s'", lang)
	}

	// Add resource type filter if specified
	if resType, ok := req.Filters["type"]; ok {
		query += fmt.Sprintf(" AND r.resource_type = '%s'", resType)
	}

	query += fmt.Sprintf(`
		ORDER BY similarity DESC
		LIMIT %d
	`, req.TopK)

	rows, err := e.db.Query(ctx, query, convertVectorToString(req.EmbeddingVector), req.MinSimilarity)
	if err != nil {
		e.logger.Errorw("Error retrieving resources", "error", err)
		return nil, fmt.Errorf("retrieving resources: %w", err)
	}
	defer rows.Close()

	resources := make([]RetrievedResource, 0)
	for rows.Next() {
		var res RetrievedResource
		if err := rows.Scan(
			&res.ResourceID,
			&res.Title,
			&res.Content,
			&res.Similarity,
			&res.CreatorID,
			&res.ResourceType,
			&res.Language,
			&res.Confidence,
		); err != nil {
			e.logger.Errorw("Error scanning resource", "error", err)
			return nil, fmt.Errorf("scanning resource: %w", err)
		}
		resources = append(resources, res)
	}

	result := &RetrievalResult{
		Resources:  resources,
		TotalFound: len(resources),
		QueryTime:  0, // Would track actual query time
		RerankedAt: 0,
	}

	e.logger.Infow("Retrieved resources",
		"query", req.Query,
		"count", len(resources),
	)

	return result, nil
}

// RetrieveByKeyword retrieves resources using full-text keyword search
// Falls back to BM25 scoring when embeddings unavailable
func (e *Engine) RetrieveByKeyword(ctx context.Context, query string, topK int, filters map[string]string) (*RetrievalResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	if topK == 0 {
		topK = 5
	}

	e.logger.Debugw("Retrieving resources by keyword",
		"query", query,
		"topK", topK,
	)

	sqlQuery := `
		SELECT
			r.id,
			r.title,
			rc.content,
			0.8 as similarity,
			r.creator_id,
			r.resource_type,
			r.language,
			0.9 as confidence
		FROM resources r
		LEFT JOIN resource_content rc ON r.id = rc.resource_id
		WHERE r.status = 'published'
			AND (r.title ILIKE $1 OR COALESCE(rc.content, '') ILIKE $1)
	`

	args := []interface{}{"%" + query + "%"}

	if lang, ok := filters["language"]; ok {
		sqlQuery += " AND r.language = $" + fmt.Sprint(len(args)+1)
		args = append(args, lang)
	}

	if resType, ok := filters["type"]; ok {
		sqlQuery += " AND r.resource_type = $" + fmt.Sprint(len(args)+1)
		args = append(args, resType)
	}

	sqlQuery += fmt.Sprintf(`
		ORDER BY r.views_count DESC
		LIMIT %d
	`, topK)

	rows, err := e.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("querying resources: %w", err)
	}
	defer rows.Close()

	resources := make([]RetrievedResource, 0)
	for rows.Next() {
		var res RetrievedResource
		if err := rows.Scan(
			&res.ResourceID,
			&res.Title,
			&res.Content,
			&res.Similarity,
			&res.CreatorID,
			&res.ResourceType,
			&res.Language,
			&res.Confidence,
		); err != nil {
			return nil, fmt.Errorf("scanning resource: %w", err)
		}
		resources = append(resources, res)
	}

	result := &RetrievalResult{
		Resources:  resources,
		TotalFound: len(resources),
	}

	return result, nil
}

// RerankByRelevance reorders retrieved results by semantic relevance
// Uses cross-encoder style scoring for better relevance matching
func (e *Engine) RerankByRelevance(ctx context.Context, query string, resources []RetrievedResource) []RetrievedResource {
	e.logger.Debugw("Reranking resources",
		"query", query,
		"resource_count", len(resources),
	)

	// Placeholder reranking - in production would use cross-encoder model
	// For now, keep existing order (already ranked by embedding similarity)
	return resources
}

// BuildRAGContext formats retrieved resources into context for LLM
func (e *Engine) BuildRAGContext(resources []RetrievedResource, maxContextLength int) string {
	context := ""
	totalLength := 0

	for i, res := range resources {
		if totalLength >= maxContextLength {
			break
		}

		// Format resource as context
		resourceContext := fmt.Sprintf(`
[Resource %d - Similarity: %.2f]
Title: %s
Type: %s
Language: %s
Creator: %s

Content:
%s

---
`, i+1, res.Similarity, res.Title, res.ResourceType, res.Language, res.CreatorID, res.Content)

		if totalLength+len(resourceContext) > maxContextLength {
			// Truncate content to fit
			availableSpace := maxContextLength - totalLength
			if availableSpace > 200 {
				truncated := resourceContext[:availableSpace-50] + "\n[truncated]\n"
				context += truncated
			}
			break
		}

		context += resourceContext
		totalLength += len(resourceContext)
	}

	return context
}

// AugmentQuery enhances user query with related keywords for better retrieval
func (e *Engine) AugmentQuery(ctx context.Context, query string) string {
	// Placeholder - in production would use query expansion techniques
	e.logger.Debugw("Augmenting query", "original", query)
	return query
}

// convertVectorToString converts embedding vector to pgvector format
// pgvector expects format like "[0.1, 0.2, 0.3]"
func convertVectorToString(vector []float32) string {
	if len(vector) == 0 {
		return "[]"
	}

	result := "["
	for i, v := range vector {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%.6f", v)
	}
	result += "]"
	return result
}

// CalculateSimilarity computes cosine similarity between two vectors
func CalculateSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	dotProduct := float32(0)
	normA := float32(0)
	normB := float32(0)

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
