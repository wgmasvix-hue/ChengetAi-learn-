package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/ai-service/internal/services"
	"go.uber.org/zap"
)

// TutorHandler handles AI tutor endpoints
type TutorHandler struct {
	aiService *services.AIService
	logger    *zap.SugaredLogger
}

// NewTutorHandler creates a new tutor handler
func NewTutorHandler(aiService *services.AIService, logger *zap.SugaredLogger) *TutorHandler {
	return &TutorHandler{
		aiService: aiService,
		logger:    logger,
	}
}

// AskQuestion handles POST /tutor/ask
func (h *TutorHandler) AskQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req services.TutorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	response, err := h.aiService.AskQuestion(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("Error answering question", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ExplainConcept handles POST /tutor/explain
func (h *TutorHandler) ExplainConcept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req services.ExplanationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	response, err := h.aiService.ExplainConcept(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("Error explaining concept", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetLearningPath handles GET /tutor/learning-path/{userId}
func (h *TutorHandler) GetLearningPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("userId")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID required"})
		return
	}

	var req services.LearningPathRequest
	req.UserID = userID

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
			return
		}
	}

	response, err := h.aiService.GenerateLearningPath(r.Context(), &req)
	if err != nil {
		h.logger.Errorw("Error generating learning path", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
