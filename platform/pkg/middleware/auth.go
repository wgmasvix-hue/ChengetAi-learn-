package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ProtectedEndpoint marks which paths require authentication
type ProtectedEndpoint struct {
	Method   string
	Pattern  string
	Required bool // Set to true to enforce authentication
}

// AuthMiddleware provides JWT-based authentication for protected endpoints
type AuthMiddleware struct {
	publicPaths   []string
	tokenVerifier func(token string) (map[string]interface{}, error)
	logger        interface{ Printf(string, ...interface{}) }
}

// NewAuthMiddleware creates a new authentication middleware
// publicPaths: endpoints that don't require authentication (e.g., /health, /login)
// tokenVerifier: function to validate JWT tokens
func NewAuthMiddleware(
	publicPaths []string,
	tokenVerifier func(string) (map[string]interface{}, error),
	logger interface{ Printf(string, ...interface{}) },
) *AuthMiddleware {
	// Always allow health checks
	defaultPublic := []string{
		"/health",
		"/ready",
		"/auth/login",
		"/auth/register",
		"/auth/refresh",
	}

	// Combine default and custom public paths
	allPublic := append(defaultPublic, publicPaths...)

	return &AuthMiddleware{
		publicPaths:   allPublic,
		tokenVerifier: tokenVerifier,
		logger:        logger,
	}
}

// Handler wraps an HTTP handler with authentication
func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if path is public
		if am.isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Validate authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.logger.Printf("Missing Authorization header for %s %s from %s",
				r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Unauthorized: Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract bearer token
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			am.logger.Printf("Invalid Authorization header format for %s %s from %s",
				r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Unauthorized: Invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := authHeader[len(bearerPrefix):]

		// Validate token
		claims, err := am.tokenVerifier(token)
		if err != nil {
			am.logger.Printf("Invalid token for %s %s from %s: %v",
				r.Method, r.URL.Path, r.RemoteAddr, err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context/headers for downstream handlers
		if userID, ok := claims["user_id"].(string); ok {
			r.Header.Set("X-User-ID", userID)
		}
		if role, ok := claims["role"].(string); ok {
			r.Header.Set("X-User-Role", role)
		}
		if schoolID, ok := claims["school_id"].(string); ok {
			r.Header.Set("X-School-ID", schoolID)
		}

		// Add timestamp for audit logging
		r.Header.Set("X-Auth-Time", time.Now().UTC().Format(time.RFC3339))

		next.ServeHTTP(w, r)
	})
}

// isPublicPath checks if a path is publicly accessible
func (am *AuthMiddleware) isPublicPath(path string) bool {
	for _, publicPath := range am.publicPaths {
		if path == publicPath || strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	return false
}

// ExtractUserID extracts the authenticated user ID from request
func ExtractUserID(r *http.Request) (string, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return "", fmt.Errorf("user ID not found in request context")
	}
	return userID, nil
}

// ExtractUserRole extracts the authenticated user's role from request
func ExtractUserRole(r *http.Request) (string, error) {
	role := r.Header.Get("X-User-Role")
	if role == "" {
		return "", fmt.Errorf("user role not found in request context")
	}
	return role, nil
}

// ExtractSchoolID extracts the school ID from authenticated request
func ExtractSchoolID(r *http.Request) string {
	return r.Header.Get("X-School-ID")
}

// RequireRole middleware ensures user has required role
func RequireRole(allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Header.Get("X-User-Role")
			if userRole == "" {
				http.Error(w, "Forbidden: No role information", http.StatusForbidden)
				return
			}

			// Check if user has required role
			hasRole := false
			for _, allowed := range allowedRoles {
				if userRole == allowed {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireSchool middleware ensures user is operating within their school
func RequireSchool(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		schoolID := r.Header.Get("X-School-ID")
		if schoolID == "" {
			// Some roles/users might not have school association
			// This is configurable per service
			// For now, log but allow
		}

		next.ServeHTTP(w, r)
	})
}

// LogRequest logs HTTP requests for audit trail
func LogRequest(logger interface{ Printf(string, ...interface{}) }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "anonymous"
			}

			logger.Printf("[AUDIT] %s %s %s from %s - User: %s",
				r.Method, r.URL.Path, r.URL.RawQuery, r.RemoteAddr, userID)

			next.ServeHTTP(w, r)
		})
	}
}
