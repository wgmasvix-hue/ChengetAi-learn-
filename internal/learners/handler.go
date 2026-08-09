package learners

import (
	"encoding/json"
	stdhttp "net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Me(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, stdhttp.StatusUnauthorized, "UNAUTHORIZED", "Authentication token is required")
		return
	}
	profile, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, stdhttp.StatusNotFound, "NOT_FOUND", "Learner profile was not found")
		return
	}
	writeJSON(w, stdhttp.StatusOK, profile)
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w stdhttp.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
