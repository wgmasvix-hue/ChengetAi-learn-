package embeddings

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Embedder generates embeddings for text
type Embedder struct {
	model  string
	logger *zap.SugaredLogger
}

// NewEmbedder creates a new embedder
func NewEmbedder(model string, logger *zap.SugaredLogger) *Embedder {
	return &Embedder{
		model:  model,
		logger: logger,
	}
}

// Embedding represents a vector embedding
type Embedding struct {
	Vector []float32
	Model  string
	Dimension int
}

// GenerateEmbedding generates an embedding for text
// This is a placeholder implementation
// In production, this would call a real embedding service
// (e.g., Ollama, OpenAI API, Hugging Face, or local model)
func (e *Embedder) GenerateEmbedding(ctx context.Context, text string) (*Embedding, error) {
	if text == "" {
		return nil, fmt.Errorf("empty text")
	}

	// Placeholder: Generate a simple embedding based on text characteristics
	// In production, replace with real embedding model call
	embedding := e.generatePlaceholderEmbedding(text)

	e.logger.Debugw("Generated embedding",
		"model", e.model,
		"dimension", len(embedding.Vector),
		"text_length", len(text),
	)

	return embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (e *Embedder) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([]Embedding, error) {
	embeddings := make([]Embedding, len(texts))

	for i, text := range texts {
		emb, err := e.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("embedding text %d: %w", i, err)
		}
		embeddings[i] = *emb
	}

	e.logger.Infow("Generated batch embeddings",
		"model", e.model,
		"count", len(embeddings),
	)

	return embeddings, nil
}

// generatePlaceholderEmbedding creates a placeholder embedding
// This should be replaced with actual embedding model in production
func (e *Embedder) generatePlaceholderEmbedding(text string) *Embedding {
	// For production:
	// - Use sentence-transformers library via Python gRPC service
	// - Or use OpenAI API embeddings endpoint
	// - Or use local Ollama instance
	// - Or use Hugging Face inference API

	// Placeholder: Create vector based on text hash
	const dimension = 384

	vector := make([]float32, dimension)

	// Simple hash-based approach for placeholder
	hash := 0
	for _, char := range text {
		hash = ((hash << 5) - hash) + int(char)
	}

	// Fill vector with pseudo-random values based on hash
	for i := 0; i < dimension; i++ {
		// Deterministic "random" values based on hash and index
		seed := uint32(hash) ^ uint32(i*7919)
		// Linear congruential generator
		seed = seed*1103515245 + 12345
		value := float32((seed/65536)%32768) / 32768.0
		vector[i] = (value * 2) - 1 // Range [-1, 1]
	}

	// Normalize the vector
	var sum float32
	for _, v := range vector {
		sum += v * v
	}

	if sum > 0 {
		norm := float32(1.0) / float32(len(vector))
		for i := range vector {
			vector[i] *= norm
		}
	}

	return &Embedding{
		Vector:    vector,
		Model:     e.model,
		Dimension: dimension,
	}
}

// CosineSimilarity calculates similarity between two embeddings
func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
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

	return dotProduct / (float32(len(a)) * (normA * normB))
}
