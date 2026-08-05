package errors

import (
	"fmt"
)

// AppError represents an application error
type AppError struct {
	Code    string // Machine-readable error code
	Message string // Human-readable error message
	Err     error  // Underlying error
	Status  int    // HTTP status code
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrInvalidRequest    = "INVALID_REQUEST"
	ErrUnauthorized      = "UNAUTHORIZED"
	ErrPermissionDenied   = "PERMISSION_DENIED"
	ErrNotFound          = "NOT_FOUND"
	ErrConflict          = "CONFLICT"
	ErrInternal          = "INTERNAL_ERROR"
	ErrEmailTaken        = "EMAIL_TAKEN"
	ErrUsernameTaken     = "USERNAME_TAKEN"
	ErrInvalidPassword   = "INVALID_PASSWORD"
	ErrSessionExpired    = "SESSION_EXPIRED"
	ErrInvalidToken      = "INVALID_TOKEN"
)

// HTTP Status Codes
const (
	StatusBadRequest       = 400
	StatusUnauthorized     = 401
	StatusForbidden        = 403
	StatusNotFound         = 404
	StatusConflict         = 409
	StatusInternalError    = 500
)

// Constructor functions
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:    ErrInvalidRequest,
		Message: message,
		Err:     err,
		Status:  StatusBadRequest,
	}
}

func NewUnauthorized(message string, err error) *AppError {
	return &AppError{
		Code:    ErrUnauthorized,
		Message: message,
		Err:     err,
		Status:  StatusUnauthorized,
	}
}

func NewForbidden(message string, err error) *AppError {
	return &AppError{
		Code:    ErrPermissionDenied,
		Message: message,
		Err:     err,
		Status:  StatusForbidden,
	}
}

func NewNotFound(message string, err error) *AppError {
	return &AppError{
		Code:    ErrNotFound,
		Message: message,
		Err:     err,
		Status:  StatusNotFound,
	}
}

func NewConflict(message string, err error) *AppError {
	return &AppError{
		Code:    ErrConflict,
		Message: message,
		Err:     err,
		Status:  StatusConflict,
	}
}

func NewInternal(message string, err error) *AppError {
	return &AppError{
		Code:    ErrInternal,
		Message: message,
		Err:     err,
		Status:  StatusInternalError,
	}
}

func NewEmailTaken() *AppError {
	return &AppError{
		Code:    ErrEmailTaken,
		Message: "Email address already in use",
		Status:  StatusConflict,
	}
}

func NewUsernameTaken() *AppError {
	return &AppError{
		Code:    ErrUsernameTaken,
		Message: "Username already in use",
		Status:  StatusConflict,
	}
}

func NewInvalidPassword() *AppError {
	return &AppError{
		Code:    ErrInvalidPassword,
		Message: "Email or password is incorrect",
		Status:  StatusUnauthorized,
	}
}

func NewSessionExpired() *AppError {
	return &AppError{
		Code:    ErrSessionExpired,
		Message: "Session has expired",
		Status:  StatusUnauthorized,
	}
}

func NewInvalidToken() *AppError {
	return &AppError{
		Code:    ErrInvalidToken,
		Message: "Invalid or malformed token",
		Status:  StatusUnauthorized,
	}
}
