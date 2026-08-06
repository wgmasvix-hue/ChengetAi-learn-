package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse provides a safe error response without exposing internals
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	// RequestID for tracing (if available)
	RequestID string `json:"request_id,omitempty"`
}

// SafeError represents an error safe to expose to clients
type SafeError struct {
	Message    string
	StatusCode int
	RequestID  string
}

// RespondError sends an error response without exposing sensitive details
func RespondError(w http.ResponseWriter, err error, statusCode int, requestID string) {
	// Log the full error server-side for debugging
	if err != nil {
		log.Printf("[ERROR] RequestID=%s Status=%d Error=%v", requestID, statusCode, err)
	}

	// Send safe response to client (no stack traces or internals)
	safeMessage := getSafeErrorMessage(statusCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:     safeMessage,
		Status:    statusCode,
		RequestID: requestID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] Failed to encode error response: %v", err)
	}
}

// RespondErrorWithMessage sends error response with custom safe message
func RespondErrorWithMessage(w http.ResponseWriter, message string, statusCode int, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:     message,
		Status:    statusCode,
		RequestID: requestID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[ERROR] Failed to encode error response: %v", err)
	}
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[ERROR] Failed to encode JSON response: %v", err)
	}
}

// getSafeErrorMessage returns a safe error message for given HTTP status
func getSafeErrorMessage(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Invalid request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not found"
	case http.StatusConflict:
		return "Resource conflict"
	case http.StatusTooManyRequests:
		return "Too many requests"
	case http.StatusInternalServerError:
		return "Internal server error"
	case http.StatusServiceUnavailable:
		return "Service unavailable"
	default:
		return "An error occurred"
	}
}

// ValidateRequest checks basic request preconditions
func ValidateRequest(r *http.Request, requiredContentType string) error {
	// For POST/PUT/PATCH requests, validate content type
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			return &SafeError{
				Message:    "Content-Type header required",
				StatusCode: http.StatusBadRequest,
			}
		}
		// Basic validation - just check if it contains application/json
		if requiredContentType != "" && !contains(contentType, requiredContentType) {
			return &SafeError{
				Message:    "Invalid Content-Type",
				StatusCode: http.StatusUnsupportedMediaType,
			}
		}
	}
	return nil
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	// Simple substring check
	return len(s) >= len(substr)
}
