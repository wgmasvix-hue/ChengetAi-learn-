package auth

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "INVALID_REQUEST", "Request body is invalid")
		return
	}

	response, err := h.service.Register(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusCreated, response)
}

func (h *Handler) Login(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "INVALID_REQUEST", "Request body is invalid")
		return
	}

	response, err := h.service.Login(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, response)
}

func (h *Handler) Logout(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	token, err := middleware.TokenFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		writeError(w, stdhttp.StatusUnauthorized, "UNAUTHORIZED", "Authentication token is required")
		return
	}
	if err := h.service.Logout(r.Context(), token); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "logged_out"})
}

func (h *Handler) Me(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, stdhttp.StatusUnauthorized, "UNAUTHORIZED", "Authentication token is required")
		return
	}
	user, err := h.service.Me(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, stdhttp.StatusOK, user)
}

func (h *Handler) handleError(w stdhttp.ResponseWriter, err error) {
	var validationErr *ValidationError
	switch {
	case errors.As(err, &validationErr):
		writeError(w, stdhttp.StatusBadRequest, "INVALID_REQUEST", validationErr.Message)
	case errors.Is(err, ErrForbiddenRole):
		writeError(w, stdhttp.StatusForbidden, "FORBIDDEN", "This role cannot self-register")
	case errors.Is(err, ErrEmailAlreadyExists):
		writeError(w, stdhttp.StatusConflict, "CONFLICT", "Email address is already registered")
	case errors.Is(err, ErrInvalidCredentials):
		writeError(w, stdhttp.StatusUnauthorized, "UNAUTHORIZED", "Email or password is incorrect")
	case errors.Is(err, ErrUserNotFound):
		writeError(w, stdhttp.StatusNotFound, "NOT_FOUND", "User was not found")
	default:
		writeError(w, stdhttp.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w stdhttp.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
