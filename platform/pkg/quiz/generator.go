package quiz

import (
	"context"
	"fmt"
	"math"
	"sort"

	"go.uber.org/zap"
)

// Generator creates educational quizzes from content using LLM
type Generator struct {
	logger *zap.SugaredLogger
}

// QuizRequest specifies parameters for quiz generation
type QuizRequest struct {
	ResourceID      string
	Content         string
	QuestionCount   int
	DifficultyLevel string // easy, medium, hard
	QuestionTypes   []string // multiple_choice, short_answer, true_false, essay
	Language        string
	GradeLevel      string
}

// Quiz represents a complete quiz
type Quiz struct {
	ID              string
	ResourceID      string
	Title           string
	Description     string
	Questions       []Question
	DifficultyLevel string
	EstimatedTime   int // minutes
	PassingScore    int // 0-100
	CreatedAt       int64
}

// Question represents a single quiz question
type Question struct {
	ID              string
	Type            string // multiple_choice, short_answer, true_false, essay
	Text            string
	Options         []QuestionOption // for multiple choice
	CorrectAnswer   string
	Explanation     string
	Points          int
	DifficultyScore float32 // 0-1
	ContentKeywords []string
}

// QuestionOption represents an answer option
type QuestionOption struct {
	ID    string
	Text  string
	Index int
}

// QuizResponse represents a submitted quiz response
type QuizResponse struct {
	ID          string
	QuizID      string
	UserID      string
	Responses   []Response
	Score       int
	Percentage  float32
	SubmittedAt int64
	TimeSpent   int // seconds
}

// Response represents an answer to a question
type Response struct {
	QuestionID string
	AnswerText string
	IsCorrect  bool
	Score      int
}

// NewGenerator creates a new quiz generator
func NewGenerator(logger *zap.SugaredLogger) *Generator {
	return &Generator{
		logger: logger,
	}
}

// GenerateQuiz creates a quiz from content
// In production, would use Claude API to generate questions
func (g *Generator) GenerateQuiz(ctx context.Context, req *QuizRequest) (*Quiz, error) {
	if req.ResourceID == "" {
		return nil, fmt.Errorf("resource ID cannot be empty")
	}

	if req.QuestionCount == 0 {
		req.QuestionCount = 10
	}

	if req.DifficultyLevel == "" {
		req.DifficultyLevel = "medium"
	}

	if req.Language == "" {
		req.Language = "en"
	}

	g.logger.Debugw("Generating quiz",
		"resource_id", req.ResourceID,
		"question_count", req.QuestionCount,
		"difficulty", req.DifficultyLevel,
	)

	// Parse content to identify key concepts
	keywords := extractKeywords(req.Content)

	// Generate questions based on content and parameters
	questions := g.generateQuestions(req, keywords)

	// Ensure we have requested number of questions
	if len(questions) > req.QuestionCount {
		questions = questions[:req.QuestionCount]
	}

	// Calculate passing score
	passingScore := 70
	if req.DifficultyLevel == "hard" {
		passingScore = 75
	}

	// Estimate time (about 2 minutes per question for medium difficulty)
	timeMultiplier := 2.0
	if req.DifficultyLevel == "easy" {
		timeMultiplier = 1.5
	} else if req.DifficultyLevel == "hard" {
		timeMultiplier = 3.0
	}
	estimatedTime := int(float32(req.QuestionCount) * float32(timeMultiplier))

	quiz := &Quiz{
		ID:              generateID(),
		ResourceID:      req.ResourceID,
		Title:           fmt.Sprintf("Quiz: %s", req.ResourceID[:20]),
		Description:     fmt.Sprintf("Assessment with %d %s questions", req.QuestionCount, req.DifficultyLevel),
		Questions:       questions,
		DifficultyLevel: req.DifficultyLevel,
		EstimatedTime:   estimatedTime,
		PassingScore:    passingScore,
		CreatedAt:       getCurrentTimestamp(),
	}

	g.logger.Infow("Generated quiz",
		"quiz_id", quiz.ID,
		"question_count", len(quiz.Questions),
	)

	return quiz, nil
}

// generateQuestions creates individual questions from content
func (g *Generator) generateQuestions(req *QuizRequest, keywords []string) []Question {
	questions := make([]Question, 0)

	// Generate multiple choice questions
	if contains(req.QuestionTypes, "multiple_choice") {
		count := int(math.Ceil(float64(req.QuestionCount) * 0.5))
		questions = append(questions, g.generateMultipleChoice(req, keywords, count)...)
	}

	// Generate true/false questions
	if contains(req.QuestionTypes, "true_false") {
		count := int(math.Ceil(float64(req.QuestionCount) * 0.3))
		questions = append(questions, g.generateTrueFalse(req, keywords, count)...)
	}

	// Generate short answer questions
	if contains(req.QuestionTypes, "short_answer") {
		count := int(math.Ceil(float64(req.QuestionCount) * 0.2))
		questions = append(questions, g.generateShortAnswer(req, keywords, count)...)
	}

	// If no specific types requested, generate mix
	if len(req.QuestionTypes) == 0 {
		questions = append(questions, g.generateMultipleChoice(req, keywords, req.QuestionCount/2)...)
		questions = append(questions, g.generateTrueFalse(req, keywords, req.QuestionCount/4)...)
		questions = append(questions, g.generateShortAnswer(req, keywords, req.QuestionCount/4)...)
	}

	return questions[:min(len(questions), req.QuestionCount)]
}

// generateMultipleChoice creates multiple choice questions
func (g *Generator) generateMultipleChoice(req *QuizRequest, keywords []string, count int) []Question {
	questions := make([]Question, 0, count)

	difficultyScore := getDifficultyScore(req.DifficultyLevel)

	for i := 0; i < count && i < len(keywords); i++ {
		options := []QuestionOption{
			{ID: "a", Text: "Correct answer", Index: 0},
			{ID: "b", Text: "Distractor 1", Index: 1},
			{ID: "c", Text: "Distractor 2", Index: 2},
			{ID: "d", Text: "Distractor 3", Index: 3},
		}

		question := Question{
			ID:              generateID(),
			Type:            "multiple_choice",
			Text:            fmt.Sprintf("Which of the following relates to %s?", keywords[i]),
			Options:         options,
			CorrectAnswer:   "a",
			Explanation:     fmt.Sprintf("%s is a key concept in this material", keywords[i]),
			Points:          10,
			DifficultyScore: difficultyScore,
			ContentKeywords: []string{keywords[i]},
		}
		questions = append(questions, question)
	}

	return questions
}

// generateTrueFalse creates true/false questions
func (g *Generator) generateTrueFalse(req *QuizRequest, keywords []string, count int) []Question {
	questions := make([]Question, 0, count)

	difficultyScore := getDifficultyScore(req.DifficultyLevel)

	for i := 0; i < count && i < len(keywords); i++ {
		isTrue := i%2 == 0
		correctAnswer := "True"
		if !isTrue {
			correctAnswer = "False"
		}

		question := Question{
			ID:              generateID(),
			Type:            "true_false",
			Text:            fmt.Sprintf("True or False: %s is a fundamental concept in this topic.", keywords[i]),
			Options:         []QuestionOption{{ID: "t", Text: "True"}, {ID: "f", Text: "False"}},
			CorrectAnswer:   correctAnswer,
			Explanation:     fmt.Sprintf("This statement is %s because...", correctAnswer),
			Points:          5,
			DifficultyScore: difficultyScore,
			ContentKeywords: []string{keywords[i]},
		}
		questions = append(questions, question)
	}

	return questions
}

// generateShortAnswer creates short answer questions
func (g *Generator) generateShortAnswer(req *QuizRequest, keywords []string, count int) []Question {
	questions := make([]Question, 0, count)

	difficultyScore := getDifficultyScore(req.DifficultyLevel)

	for i := 0; i < count && i < len(keywords); i++ {
		question := Question{
			ID:              generateID(),
			Type:            "short_answer",
			Text:            fmt.Sprintf("Briefly explain what %s means in this context.", keywords[i]),
			CorrectAnswer:   fmt.Sprintf("Definition of %s", keywords[i]),
			Explanation:     fmt.Sprintf("%s refers to a specific aspect of the content", keywords[i]),
			Points:          15,
			DifficultyScore: difficultyScore,
			ContentKeywords: []string{keywords[i]},
		}
		questions = append(questions, question)
	}

	return questions
}

// ScoreResponse evaluates a submitted quiz response
func (g *Generator) ScoreResponse(quiz *Quiz, response *QuizResponse) (*QuizResponse, error) {
	if len(quiz.Questions) != len(response.Responses) {
		g.logger.Warnw("Response count mismatch",
			"expected", len(quiz.Questions),
			"received", len(response.Responses),
		)
	}

	totalScore := 0
	totalPoints := 0

	for _, resp := range response.Responses {
		// Find corresponding question
		var question *Question
		for _, q := range quiz.Questions {
			if q.ID == resp.QuestionID {
				question = &q
				break
			}
		}

		if question == nil {
			continue
		}

		totalPoints += question.Points

		// Score the response
		if isCorrect(resp.AnswerText, question) {
			resp.IsCorrect = true
			resp.Score = question.Points
			totalScore += question.Points
		}
	}

	if totalPoints > 0 {
		response.Percentage = float32(totalScore*100) / float32(totalPoints)
	}
	response.Score = totalScore

	g.logger.Infow("Scored quiz response",
		"quiz_id", quiz.ID,
		"score", response.Score,
		"percentage", response.Percentage,
	)

	return response, nil
}

// extractKeywords extracts key concepts from content
func extractKeywords(content string) []string {
	// Placeholder - in production would use NLP or LLM
	keywords := []string{
		"concept_1", "concept_2", "concept_3",
		"concept_4", "concept_5", "concept_6",
		"concept_7", "concept_8", "concept_9",
		"concept_10",
	}
	return keywords
}

// isCorrect checks if response matches correct answer
func isCorrect(answer string, question *Question) bool {
	// Placeholder - in production would use semantic similarity for short answers
	return answer == question.CorrectAnswer
}

// getDifficultyScore converts difficulty level to score
func getDifficultyScore(level string) float32 {
	switch level {
	case "easy":
		return 0.3
	case "medium":
		return 0.6
	case "hard":
		return 0.9
	default:
		return 0.5
	}
}

// generateID creates unique identifier
func generateID() string {
	return fmt.Sprintf("id-%d", getCurrentTimestamp())
}

// getCurrentTimestamp returns current Unix timestamp
func getCurrentTimestamp() int64 {
	// Placeholder
	return 1704067200 // Example timestamp
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// contains checks if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetDifficultyDistribution returns recommended distribution of question difficulties
func GetDifficultyDistribution(totalQuestions int, level string) map[string]int {
	distribution := make(map[string]int)

	switch level {
	case "easy":
		distribution["easy"] = int(float32(totalQuestions) * 0.6)
		distribution["medium"] = int(float32(totalQuestions) * 0.3)
		distribution["hard"] = int(float32(totalQuestions) * 0.1)
	case "medium":
		distribution["easy"] = int(float32(totalQuestions) * 0.2)
		distribution["medium"] = int(float32(totalQuestions) * 0.6)
		distribution["hard"] = int(float32(totalQuestions) * 0.2)
	case "hard":
		distribution["easy"] = int(float32(totalQuestions) * 0.1)
		distribution["medium"] = int(float32(totalQuestions) * 0.3)
		distribution["hard"] = int(float32(totalQuestions) * 0.6)
	default:
		distribution["medium"] = totalQuestions
	}

	return distribution
}
