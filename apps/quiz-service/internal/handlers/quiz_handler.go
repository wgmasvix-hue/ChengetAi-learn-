package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/quiz-service/internal/services"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/quiz"
	"go.uber.org/zap"
)

// QuizHandler handles quiz endpoints
type QuizHandler struct {
	quizService *services.QuizService
	logger      *zap.SugaredLogger
}

// NewQuizHandler creates a new quiz handler
func NewQuizHandler(quizService *services.QuizService, logger *zap.SugaredLogger) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		logger:      logger,
	}
}

// GenerateQuiz handles POST /quizzes/generate
func (h *QuizHandler) GenerateQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req quiz.QuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	q, err := h.quizService.GenerateQuiz(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("Error generating quiz", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

// GetQuiz handles GET /quizzes/{quizId}
func (h *QuizHandler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	quizID := r.PathValue("quizId")
	if quizID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Quiz ID required"})
		return
	}

	q, err := h.quizService.GetQuiz(r.Context(), quizID)
	if err != nil {
		h.logger.Errorw("Error retrieving quiz", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(q)
}

// SubmitResponse handles POST /quizzes/{quizId}/submit
func (h *QuizHandler) SubmitResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	quizID := r.PathValue("quizId")
	if quizID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Quiz ID required"})
		return
	}

	var req struct {
		UserID    string          `json:"userId"`
		Responses []quiz.Response `json:"responses"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}

	response, err := h.quizService.SubmitResponse(r.Context(), quizID, req.UserID, req.Responses)
	if err != nil {
		h.logger.Errorw("Error submitting response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetUserQuizzes handles GET /quizzes/user/{userId}
func (h *QuizHandler) GetUserQuizzes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("userId")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		_, _ = sscanf(l, "%d", &limit)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		_, _ = sscanf(o, "%d", &offset)
	}

	quizzes, err := h.quizService.GetUserQuizzes(r.Context(), userID, limit, offset)
	if err != nil {
		h.logger.Errorw("Error retrieving user quizzes", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quizzes)
}

// GetResponse handles GET /responses/{responseId}
func (h *QuizHandler) GetResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	responseID := r.PathValue("responseId")
	if responseID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Response ID required"})
		return
	}

	response, err := h.quizService.GetResponse(r.Context(), responseID)
	if err != nil {
		h.logger.Errorw("Error retrieving response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function
func sscanf(input, format string, a ...interface{}) (int, error) {
	_, err := json.Unmarshal([]byte(input), a[0])
	if err == nil {
		return 1, nil
	}
	return 0, err
}
