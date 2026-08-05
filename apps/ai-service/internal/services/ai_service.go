package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/llm"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/rag"
	"go.uber.org/zap"
)

// AIService provides AI-powered educational features using RAG and LLM
type AIService struct {
	claude    *llm.ClaudeClient
	rag       *rag.Engine
	db        *pgxpool.Pool
	logger    *zap.SugaredLogger
}

// TutorRequest represents a request to the AI tutor
type TutorRequest struct {
	UserID   string
	Query    string
	Context  string // Optional context (e.g., resource ID)
	Language string
}

// TutorResponse represents AI tutor's response
type TutorResponse struct {
	Answer       string
	Sources      []string // Resource IDs used
	Explanation  string
	FollowUpQuestions []string
	Confidence   float32
	GeneratedAt  int64
}

// ExplanationRequest requests explanation of a concept
type ExplanationRequest struct {
	Concept      string
	ResourceID   string
	DetailLevel  string // brief, standard, detailed
	Language     string
}

// LearningPathRequest requests personalized learning path
type LearningPathRequest struct {
	UserID      string
	Goal        string
	CurrentLevel string // beginner, intermediate, advanced
	TimeAvailable int    // hours per week
}

// LearningPath represents a personalized learning plan
type LearningPath struct {
	UserID        string
	Goal          string
	CurrentLevel  string
	Milestones    []Milestone
	EstimatedTime int // weeks
	Resources     []string // Resource IDs
	CreatedAt     int64
}

// Milestone represents a learning milestone
type Milestone struct {
	Title       string
	Description string
	Resources   []string
	Assessments []string
	Order       int
}

// NewAIService creates a new AI service
func NewAIService(claude *llm.ClaudeClient, ragEngine *rag.Engine, db *pgxpool.Pool, logger *zap.SugaredLogger) *AIService {
	return &AIService{
		claude: claude,
		rag:    ragEngine,
		db:     db,
		logger: logger,
	}
}

// AskQuestion handles general questions from the AI tutor
func (s *AIService) AskQuestion(ctx context.Context, req *TutorRequest) (*TutorResponse, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	if req.Language == "" {
		req.Language = "en"
	}

	s.logger.Debugw("AI tutor question received",
		"user_id", req.UserID,
		"query", req.Query,
	)

	// Step 1: Retrieve relevant context using RAG
	retrievalReq := &rag.RetrievalRequest{
		Query:          req.Query,
		TopK:           5,
		MinSimilarity:  0.5,
		Filters:        map[string]string{"language": req.Language},
	}

	// In production, would use actual embedding from query
	placeholderEmbedding := make([]float32, 384)
	for i := range placeholderEmbedding {
		placeholderEmbedding[i] = 0.1
	}
	retrievalReq.EmbeddingVector = placeholderEmbedding

	retrievalResult, err := s.rag.RetrieveByEmbedding(ctx, retrievalReq)
	if err != nil {
		s.logger.Errorw("Error retrieving context", "error", err)
		// Continue without retrieval if it fails
		retrievalResult = &rag.RetrievalResult{Resources: make([]rag.RetrievedResource, 0)}
	}

	// Step 2: Build context from retrieved resources
	ragContext := s.rag.BuildRAGContext(retrievalResult.Resources, 8192)

	// Step 3: Generate system prompt for educational context
	systemPrompt := fmt.Sprintf(`You are an expert educational AI tutor powered by ChengetAi.
Your role is to:
1. Answer questions using provided resources from African educational databases
2. Explain concepts clearly at the appropriate level
3. Provide examples relevant to the African context
4. Encourage critical thinking
5. Suggest related topics for further learning

Use the provided resources to ground your answers. If resources don't contain relevant information, say so clearly.
Always be encouraging and supportive in your teaching approach.`)

	// Step 4: Build and send prompt to Claude
	completionReq := s.claude.BuildRAGPrompt(req.Query,
		resourcesListToStrings(retrievalResult.Resources),
		systemPrompt)

	completionRes, err := s.claude.GenerateCompletion(ctx, completionReq)
	if err != nil {
		s.logger.Errorw("Error generating completion", "error", err)
		return nil, fmt.Errorf("generating response: %w", err)
	}

	// Step 5: Extract sources and generate follow-up questions
	sourceIDs := extractSourceIDs(retrievalResult.Resources)
	followUps := generateFollowUpQuestions(req.Query, completionRes.Content)

	response := &TutorResponse{
		Answer:       completionRes.Content,
		Sources:      sourceIDs,
		Explanation:  ragContext,
		FollowUpQuestions: followUps,
		Confidence:   0.85, // Placeholder
		GeneratedAt:  getCurrentTimestamp(),
	}

	s.logger.Infow("Generated tutor response",
		"user_id", req.UserID,
		"source_count", len(sourceIDs),
	)

	return response, nil
}

// ExplainConcept generates detailed explanation of a concept
func (s *AIService) ExplainConcept(ctx context.Context, req *ExplanationRequest) (*TutorResponse, error) {
	if req.Concept == "" {
		return nil, fmt.Errorf("concept cannot be empty")
	}

	if req.Language == "" {
		req.Language = "en"
	}

	if req.DetailLevel == "" {
		req.DetailLevel = "standard"
	}

	s.logger.Debugw("Concept explanation requested",
		"concept", req.Concept,
		"detail_level", req.DetailLevel,
	)

	// Retrieve resources about the concept
	retrievalReq := &rag.RetrievalRequest{
		Query:          req.Concept,
		TopK:           5,
		MinSimilarity:  0.4,
		Filters:        map[string]string{"language": req.Language},
	}

	placeholderEmbedding := make([]float32, 384)
	for i := range placeholderEmbedding {
		placeholderEmbedding[i] = 0.1
	}
	retrievalReq.EmbeddingVector = placeholderEmbedding

	retrievalResult, err := s.rag.RetrieveByEmbedding(ctx, retrievalReq)
	if err != nil {
		retrievalResult = &rag.RetrievalResult{Resources: make([]rag.RetrievedResource, 0)}
	}

	// Build detailed explanation prompt
	detailInstruction := ""
	switch req.DetailLevel {
	case "brief":
		detailInstruction = "Provide a concise explanation in 2-3 sentences."
	case "detailed":
		detailInstruction = "Provide a comprehensive explanation with examples, history, and applications."
	default:
		detailInstruction = "Provide a balanced explanation with key points and one example."
	}

	systemPrompt := fmt.Sprintf(`You are an expert science and mathematics educator.
Explain the following concept: %s

%s

Include:
- Clear definition
- How it works
- Real-world applications in African context
- Common misconceptions
- Key resources for further learning`, req.Concept, detailInstruction)

	// Generate explanation
	completionReq := &llm.CompletionRequest{
		SystemPrompt: systemPrompt,
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Please explain %s in detail.", req.Concept),
			},
		},
		Temperature: 0.7,
		MaxTokens:   2048,
	}

	completionRes, err := s.claude.GenerateCompletion(ctx, completionReq)
	if err != nil {
		return nil, fmt.Errorf("generating explanation: %w", err)
	}

	sourceIDs := extractSourceIDs(retrievalResult.Resources)

	response := &TutorResponse{
		Answer:      completionRes.Content,
		Sources:     sourceIDs,
		Confidence:  0.88,
		GeneratedAt: getCurrentTimestamp(),
	}

	return response, nil
}

// GenerateLearningPath creates a personalized learning path for a user
func (s *AIService) GenerateLearningPath(ctx context.Context, req *LearningPathRequest) (*LearningPath, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	if req.Goal == "" {
		return nil, fmt.Errorf("goal cannot be empty")
	}

	if req.CurrentLevel == "" {
		req.CurrentLevel = "beginner"
	}

	s.logger.Debugw("Learning path requested",
		"user_id", req.UserID,
		"goal", req.Goal,
		"current_level", req.CurrentLevel,
	)

	// Placeholder: Generate learning path structure
	milestones := []Milestone{
		{
			Title:       "Foundation",
			Description: "Master the basics",
			Order:       1,
			Resources:   []string{"resource-1", "resource-2"},
			Assessments: []string{"quiz-1"},
		},
		{
			Title:       "Development",
			Description: "Build intermediate skills",
			Order:       2,
			Resources:   []string{"resource-3", "resource-4"},
			Assessments: []string{"quiz-2"},
		},
		{
			Title:       "Mastery",
			Description: "Advanced application",
			Order:       3,
			Resources:   []string{"resource-5", "resource-6"},
			Assessments: []string{"quiz-3", "project-1"},
		},
	}

	estimatedWeeks := 12
	if req.TimeAvailable > 0 {
		estimatedWeeks = 36 / req.TimeAvailable
	}

	path := &LearningPath{
		UserID:        req.UserID,
		Goal:          req.Goal,
		CurrentLevel:  req.CurrentLevel,
		Milestones:    milestones,
		EstimatedTime: estimatedWeeks,
		Resources:     []string{"resource-1", "resource-2", "resource-3", "resource-4", "resource-5", "resource-6"},
		CreatedAt:     getCurrentTimestamp(),
	}

	s.logger.Infow("Generated learning path",
		"user_id", req.UserID,
		"milestones", len(milestones),
	)

	return path, nil
}

// Helper functions

func resourcesListToStrings(resources []rag.RetrievedResource) []string {
	results := make([]string, len(resources))
	for i, res := range resources {
		results[i] = fmt.Sprintf("Resource: %s\n%s", res.Title, res.Content)
	}
	return results
}

func extractSourceIDs(resources []rag.RetrievedResource) []string {
	ids := make([]string, len(resources))
	for i, res := range resources {
		ids[i] = res.ResourceID
	}
	return ids
}

func generateFollowUpQuestions(originalQuery, response string) []string {
	// Placeholder - in production would use LLM or keyword extraction
	return []string{
		"Can you give more examples?",
		"How does this relate to other concepts?",
		"What are the practical applications?",
	}
}

func getCurrentTimestamp() int64 {
	return 1704067200
}
