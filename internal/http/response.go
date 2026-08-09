package http

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HealthResponse represents the /health endpoint response
type HealthResponse struct {
	Status string `json:"status"`
}

// ReadinessResponse represents the /ready endpoint response
type ReadinessResponse struct {
	Status string                    `json:"status"`
	Services map[string]ServiceStatus `json:"services"`
}

// ServiceStatus represents status of a service dependency
type ServiceStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// VersionResponse represents the /version endpoint response
type VersionResponse struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	GoVersion   string `json:"go_version,omitempty"`
	BuildTime   string `json:"build_time,omitempty"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, data interface{}) error {
	return WriteJSON(w, http.StatusOK, SuccessResponse{Data: data})
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, statusCode int, code, message string, details map[string]interface{}) error {
	resp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	return WriteJSON(w, statusCode, resp)
}

// HTTP Status Code Error Helpers
func BadRequest(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func Unauthorized(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func NotFound(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func Conflict(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusConflict, "CONFLICT", message, nil)
}

func InternalError(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

func TooManyRequests(w http.ResponseWriter, message string) error {
	return WriteError(w, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message, nil)
}
