package knowledge

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

// Processor handles document processing for resources
type Processor struct {
	logger *zap.SugaredLogger
}

// NewProcessor creates a new knowledge processor
func NewProcessor(logger *zap.SugaredLogger) *Processor {
	return &Processor{
		logger: logger,
	}
}

// ProcessedContent represents processed content from a document
type ProcessedContent struct {
	RawText          string
	TextLength       int
	Chunks           []string
	Keywords         []string
	Summary          string
	LanguageDetected string
	ProcessedAt      int64 // timestamp
}

// ProcessDocument processes a document's content
func (p *Processor) ProcessDocument(ctx context.Context, content string, mimeType string) (*ProcessedContent, error) {
	if content == "" {
		return nil, fmt.Errorf("empty content")
	}

	processed := &ProcessedContent{
		RawText:    content,
		TextLength: len(content),
	}

	// Chunk content for embedding
	processed.Chunks = p.chunkText(content)

	// Extract keywords
	processed.Keywords = p.extractKeywords(content)

	// Detect language (simplified - in production use a proper language detection library)
	processed.LanguageDetected = p.detectLanguage(content)

	p.logger.Infow("Document processed",
		"text_length", processed.TextLength,
		"chunks", len(processed.Chunks),
		"keywords", len(processed.Keywords),
		"language", processed.LanguageDetected,
	)

	return processed, nil
}

// chunkText splits content into chunks for embedding
func (p *Processor) chunkText(text string) []string {
	const chunkSize = 512 // Characters per chunk
	const overlapSize = 64 // Overlap between chunks

	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += (chunkSize - overlapSize) {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := string(runes[i:end])
		if strings.TrimSpace(chunk) != "" {
			chunks = append(chunks, chunk)
		}

		if end == len(runes) {
			break
		}
	}

	return chunks
}

// extractKeywords extracts keywords from text (simplified)
func (p *Processor) extractKeywords(text string) []string {
	// This is a simplified implementation
	// In production, use NLP library like go-nlp or spaCy via gRPC

	words := strings.Fields(strings.ToLower(text))
	wordFreq := make(map[string]int)

	// Common stop words to filter
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "as": true, "it": true, "this": true,
		"that": true, "these": true, "those": true, "which": true, "who": true,
	}

	// Count word frequency
	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:")

		if len(word) > 3 && !stopWords[word] {
			wordFreq[word]++
		}
	}

	// Sort by frequency (simplified - just get top ones)
	var keywords []string
	for word, freq := range wordFreq {
		if freq >= 2 && len(keywords) < 20 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// detectLanguage detects the language of text (simplified)
func (p *Processor) detectLanguage(text string) string {
	// This is a simplified implementation
	// In production, use a proper language detection library

	// Count common words per language
	languages := map[string][]string{
		"en": {"the", "and", "is", "to", "a", "of", "in", "that"},
		"sw": {"na", "ni", "wa", "la", "kwa", "ya", "tu", "kuwa"},
		"zu": {"i", "u", "e", "uku", "ukuba", "okusho", "ngoku"},
		"xh": {"i", "u", "e", "uku", "ukuya", "igama", "into"},
		"fr": {"le", "de", "et", "un", "une", "la", "les", "des"},
	}

	textLower := strings.ToLower(text)
	counts := make(map[string]int)

	for lang, words := range languages {
		for _, word := range words {
			if strings.Contains(textLower, word) {
				counts[lang]++
			}
		}
	}

	// Find language with highest count
	maxCount := 0
	detected := "en"
	for lang, count := range counts {
		if count > maxCount {
			maxCount = count
			detected = lang
		}
	}

	return detected
}
