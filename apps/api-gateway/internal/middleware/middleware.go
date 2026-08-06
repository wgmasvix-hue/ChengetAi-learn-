package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/api-gateway/internal/config"
	"go.uber.org/zap"
)

type contextKey string

const (
	RequestIDKey contextKey = "request-id"
)

// Chain applies middleware in order
func Chain(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

// RequestID adds a request ID to each request
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r)
		})
	}
}

// RequestLogger logs incoming requests
func RequestLogger(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.RequestURI,
				"duration_ms", duration.Milliseconds(),
				"request_id", r.Header.Get("X-Request-ID"),
			)
		})
	}
}

// Recovery recovers from panics
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal server error"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// CORS sets up CORS headers with proper origin validation
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// SECURITY: Never use wildcard "*" in production CORS
			// Only allow requests from specified origins
			allowedOrigins := cfg.CORSAllowedOrigins // Load from environment
			isOriginAllowed := false

			// Check if origin is in allowlist
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					isOriginAllowed = true
					break
				}
			}

			// Only set CORS headers if origin is allowed
			if isOriginAllowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
				w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-Total-Count")
			} else if origin != "" {
				// Request from disallowed origin - don't set CORS headers
				// This will cause CORS check to fail in browser
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
