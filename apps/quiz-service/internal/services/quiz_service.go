package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/quiz-service/internal/repositories"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/quiz"
	"go.uber.org/zap"
)

// QuizService orchestrates quiz operations
type QuizService struct {
	quizRepo     *repositories.QuizRepository
	responseRepo *repositories.ResponseRepository
	generator    *quiz.Generator
	logger       *zap.SugaredLogger
}

// NewQuizService creates a new quiz service
func NewQuizService(
	quizRepo *repositories.QuizRepository,
	responseRepo *repositories.ResponseRepository,
	generator *quiz.Generator,
	logger *zap.SugaredLogger,
) *QuizService {
	return &QuizService{
		quizRepo:     quizRepo,
		responseRepo: responseRepo,
		generator:    generator,
		logger:       logger,
	}
}

// GenerateQuiz creates a new quiz from resource content
func (s *QuizService) GenerateQuiz(ctx context.Context, req *quiz.QuizRequest) (*quiz.Quiz, error) {
	if req.ResourceID == "" {
		return nil, fmt.Errorf("resource ID cannot be empty")
	}

	if req.Content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}

	s.logger.Debugw("Generating quiz",
		"resource_id", req.ResourceID,
		"question_count", req.QuestionCount,
	)

	// Generate quiz using quiz generator
	q, err := s.generator.GenerateQuiz(ctx, req)
	if err != nil {
		s.logger.Errorw("Error generating quiz", "error", err)
		return nil, err
	}

	// Save quiz to database
	if err := s.quizRepo.Create(ctx, q); err != nil {
		s.logger.Errorw("Error saving quiz", "error", err)
		return nil, err
	}

	s.logger.Infow("Quiz generated successfully",
		"quiz_id", q.ID,
		"question_count", len(q.Questions),
	)

	return q, nil
}

// GetQuiz retrieves a quiz by ID
func (s *QuizService) GetQuiz(ctx context.Context, quizID string) (*quiz.Quiz, error) {
	q, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		s.logger.Errorw("Error retrieving quiz", "error", err)
		return nil, err
	}

	return q, nil
}

// GetQuizzesByResource retrieves quizzes for a resource
func (s *QuizService) GetQuizzesByResource(ctx context.Context, resourceID string) ([]quiz.Quiz, error) {
	quizzes, err := s.quizRepo.GetByResourceID(ctx, resourceID)
	if err != nil {
		s.logger.Errorw("Error retrieving quizzes", "error", err)
		return nil, err
	}

	return quizzes, nil
}

// SubmitResponse processes a quiz response
func (s *QuizService) SubmitResponse(ctx context.Context, quizID, userID string, responses []quiz.Response) (*quiz.QuizResponse, error) {
	// Retrieve the quiz
	q, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		return nil, fmt.Errorf("retrieving quiz: %w", err)
	}

	// Create response object
	response := &quiz.QuizResponse{
		ID:        generateID(),
		QuizID:    quizID,
		UserID:    userID,
		Responses: responses,
	}

	// Score the response
	response, err = s.generator.ScoreResponse(q, response)
	if err != nil {
		return nil, fmt.Errorf("scoring response: %w", err)
	}

	// Record timestamp
	response.SubmittedAt = getCurrentTimestamp()

	// Save response to database
	if err := s.responseRepo.Create(ctx, response); err != nil {
		s.logger.Errorw("Error saving response", "error", err)
		return nil, err
	}

	s.logger.Infow("Quiz response submitted",
		"quiz_id", quizID,
		"user_id", userID,
		"score", response.Score,
		"percentage", response.Percentage,
	)

	return response, nil
}

// GetUserQuizzes retrieves quizzes attempted by a user
func (s *QuizService) GetUserQuizzes(ctx context.Context, userID string, limit, offset int) ([]quiz.QuizResponse, error) {
	responses, err := s.responseRepo.GetByUser(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Errorw("Error retrieving user quizzes", "error", err)
		return nil, err
	}

	return responses, nil
}

// GetResponse retrieves a specific quiz response
func (s *QuizService) GetResponse(ctx context.Context, responseID string) (*quiz.QuizResponse, error) {
	response, err := s.responseRepo.GetByID(ctx, responseID)
	if err != nil {
		s.logger.Errorw("Error retrieving response", "error", err)
		return nil, err
	}

	return response, nil
}

// GetUserStats retrieves statistics for a user's quiz attempts
func (s *QuizService) GetUserStats(ctx context.Context, userID, quizID string) (*repositories.QuizStats, error) {
	stats, err := s.responseRepo.GetUserStats(ctx, userID, quizID)
	if err != nil {
		s.logger.Errorw("Error retrieving stats", "error", err)
		return nil, err
	}

	return stats, nil
}

// DeleteQuiz removes a quiz
func (s *QuizService) DeleteQuiz(ctx context.Context, quizID string) error {
	if err := s.quizRepo.Delete(ctx, quizID); err != nil {
		s.logger.Errorw("Error deleting quiz", "error", err)
		return err
	}

	s.logger.Infow("Quiz deleted", "quiz_id", quizID)
	return nil
}

// Helper functions

func generateID() string {
	return fmt.Sprintf("resp-%d", getCurrentTimestamp())
}

func getCurrentTimestamp() int64 {
	return 1704067200
}
