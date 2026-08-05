package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/recommendations-service/internal/services"
	"go.uber.org/zap"
)

// RecommendationHandler handles recommendation endpoints
type RecommendationHandler struct {
	recService *services.RecommendationService
	logger     *zap.SugaredLogger
}

// NewRecommendationHandler creates a new recommendation handler
func NewRecommendationHandler(recService *services.RecommendationService, logger *zap.SugaredLogger) *RecommendationHandler {
	return &RecommendationHandler{
		recService: recService,
		logger:     logger,
	}
}

// GetRecommendations handles GET /recommendations/{userId}
func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
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

	topK := 10
	if k := r.URL.Query().Get("topK"); k != "" {
		if val, err := strconv.Atoi(k); err == nil {
			topK = val
		}
	}

	strategy := r.URL.Query().Get("strategy")
	if strategy == "" {
		strategy = "hybrid"
	}

	result, err := h.recService.GetRecommendations(r.Context(), userID, topK, strategy)
	if err != nil {
		h.logger.Errorw("Error generating recommendations", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GenerateRecommendations handles POST /recommendations/generate
func (h *RecommendationHandler) GenerateRecommendations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID   string `json:"userId"`
		TopK     int    `json:"topK"`
		Strategy string `json:"strategy"`
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

	if req.TopK == 0 {
		req.TopK = 10
	}

	result, err := h.recService.GetRecommendations(r.Context(), req.UserID, req.TopK, req.Strategy)
	if err != nil {
		h.logger.Errorw("Error generating recommendations", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetMetrics handles GET /metrics/{userId}
func (h *RecommendationHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
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

	metrics := h.recService.GetMetrics(r.Context(), userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// SubmitFeedback handles POST /feedback
func (h *RecommendationHandler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID     string `json:"userId"`
		ResourceID string `json:"resourceId"`
		Helpful    bool   `json:"helpful"`
		Reason     string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.UserID == "" || req.ResourceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID and Resource ID required"})
		return
	}

	if err := h.recService.SubmitFeedback(r.Context(), req.UserID, req.ResourceID, req.Helpful, req.Reason); err != nil {
		h.logger.Errorw("Error saving feedback", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "feedback recorded"})
}
