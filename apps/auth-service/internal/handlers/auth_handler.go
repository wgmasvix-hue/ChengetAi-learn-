package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/auth-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/auth-service/internal/services"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
	"go.uber.org/zap"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.SugaredLogger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, errors.NewBadRequest("invalid request body", err))
		return
	}

	// Validate request
	if req.Email == "" || req.Username == "" || req.Password == "" || req.FullName == "" {
		h.writeErrorResponse(w, errors.NewBadRequest("missing required fields", nil))
		return
	}

	// Register user
	response, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, http.StatusCreated, response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, errors.NewBadRequest("invalid request body", err))
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		h.writeErrorResponse(w, errors.NewBadRequest("email and password required", nil))
		return
	}

	// Login user
	response, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, errors.NewBadRequest("invalid request body", err))
		return
	}

	response, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, http.StatusOK, response)
}

// Response types
type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type errorResponse struct {
	Success bool         `json:"success"`
	Error   errorDetails `json:"error"`
}

type errorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeSuccessResponse writes a success response
func (h *AuthHandler) writeSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(successResponse{
		Success: true,
		Data:    data,
	})
}

// writeErrorResponse writes an error response
func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.NewInternal("internal error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	json.NewEncoder(w).Encode(errorResponse{
		Success: false,
		Error: errorDetails{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})

	h.logger.Errorw("request error",
		"code", appErr.Code,
		"message", appErr.Message,
		"status", appErr.Status,
	)
}
