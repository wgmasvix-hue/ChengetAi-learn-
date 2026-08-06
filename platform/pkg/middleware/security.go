package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// SecurityHeadersMiddleware adds important security headers to all responses
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")

		// Enable XSS protection in older browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (HSTS) - adjust max-age for production
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Reduce information leakage
		w.Header().Set("Server", "ChengetAi/1.0")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles CORS with security in mind
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			isAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed || allowed == "*" {
					isAllowed = true
					break
				}
			}

			if isAllowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware provides basic rate limiting by IP address
type RateLimiter interface {
	Allow(key string) bool
}

func RateLimitMiddlewareFunc(limiter RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP (consider X-Forwarded-For if behind proxy)
			clientIP := getClientIP(r)

			if !limiter.Allow(clientIP) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For (set by reverse proxy)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Take the first IP if multiple
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Use remote address
	if r.RemoteAddr != "" {
		// Remove port if present
		ip := strings.Split(r.RemoteAddr, ":")[0]
		return ip
	}

	return "unknown"
}

// AuthenticationMiddleware validates JWT tokens
type AuthenticationMiddleware struct {
	tokenValidator func(string) (map[string]interface{}, error)
	excludePaths   []string
}

// NewAuthenticationMiddleware creates a new authentication middleware
func NewAuthenticationMiddleware(tokenValidator func(string) (map[string]interface{}, error), excludePaths []string) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		tokenValidator: tokenValidator,
		excludePaths:   excludePaths,
	}
}

// Middleware returns the HTTP middleware function
func (a *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for excluded paths
		for _, path := range a.excludePaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract Bearer token
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := authHeader[len(bearerPrefix):]

		// Validate token
		claims, err := a.tokenValidator(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Store claims in context for downstream handlers
		r.Header.Set("X-User-ID", fmt.Sprintf("%v", claims["user_id"]))
		r.Header.Set("X-User-Role", fmt.Sprintf("%v", claims["role"]))

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests with security awareness
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't log sensitive headers
		sensitiveHeaders := []string{"Authorization", "Cookie", "X-API-Key"}
		for _, header := range sensitiveHeaders {
			if r.Header.Get(header) != "" {
				r.Header.Set(header, "[REDACTED]")
			}
		}

		next.ServeHTTP(w, r)
	})
}
