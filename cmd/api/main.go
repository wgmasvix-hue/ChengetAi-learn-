package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/wgmasvix-hue/chengetai-learn/internal/auth"
	"github.com/wgmasvix-hue/chengetai-learn/internal/config"
	"github.com/wgmasvix-hue/chengetai-learn/internal/database"
	"github.com/wgmasvix-hue/chengetai-learn/internal/health"
	httputil "github.com/wgmasvix-hue/chengetai-learn/internal/http"
	"github.com/wgmasvix-hue/chengetai-learn/internal/middleware"
	redisclient "github.com/wgmasvix-hue/chengetai-learn/internal/redis"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redis, err := redisclient.Connect(context.Background(), cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Initialize token manager
	tokenManager, err := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTExpiry, "chengetai")
	if err != nil {
		log.Fatalf("Failed to create token manager: %v", err)
	}

	// Create HTTP router
	router := http.NewServeMux()

	// Setup middleware stack
	var handler http.Handler = router
	handler = middleware.SecurityHeadersMiddleware(handler)
	handler = middleware.CORSMiddleware(strings.Split(cfg.CORSAllowedOrigins, ","))(handler)
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.RequestIDMiddleware(handler)

	// Setup health check handlers
	healthHandler := health.NewHandler(db, redis)
	router.HandleFunc("GET /health", healthHandler.Health)
	router.HandleFunc("GET /ready", healthHandler.Ready)
	router.HandleFunc("GET /version", healthHandler.Version)

	// Setup API routes
	setupAuthRoutes(router, db, redis, tokenManager)
	setupUserRoutes(router, db)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	log.Printf("Starting ChengetAi Learn API on %s", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// setupAuthRoutes sets up authentication routes
func setupAuthRoutes(router *http.ServeMux, db *database.Connection, redis *redisclient.Client, tm *auth.TokenManager) {
	router.HandleFunc("POST /api/v1/auth/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, db, tm)
	})
	router.HandleFunc("POST /api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, db, tm)
	})
	router.HandleFunc("GET /api/v1/auth/me", func(w http.ResponseWriter, r *http.Request) {
		handleAuthMe(w, r, db, tm)
	})
}

// setupUserRoutes sets up user routes
func setupUserRoutes(router *http.ServeMux, db *database.Connection) {
	router.HandleFunc("GET /api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		httputil.WriteSuccess(w, map[string]string{"message": "Users endpoint"})
	})
}

// Auth handlers (minimal implementation for Phase One)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type AuthResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}

func handleRegister(w http.ResponseWriter, r *http.Request, db *database.Connection, tm *auth.TokenManager) {
	if r.Method != http.MethodPost {
		httputil.BadRequest(w, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := parseJSON(r, &req); err != nil {
		httputil.BadRequest(w, "Invalid request body")
		return
	}

	// Validate
	if req.Email == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" {
		httputil.BadRequest(w, "Missing required fields")
		return
	}

	// For Phase One, we'll just return a placeholder
	// Real implementation will insert into database
	httputil.WriteJSON(w, http.StatusCreated, AuthResponse{
		UserID: "550e8400-e29b-41d4-a716-446655440000",
		Email:  req.Email,
		Role:   req.Role,
		Token:  "placeholder-token",
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request, db *database.Connection, tm *auth.TokenManager) {
	if r.Method != http.MethodPost {
		httputil.BadRequest(w, "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := parseJSON(r, &req); err != nil {
		httputil.BadRequest(w, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		httputil.BadRequest(w, "Email and password required")
		return
	}

	// For Phase One, placeholder
	httputil.WriteJSON(w, http.StatusOK, AuthResponse{
		UserID: "550e8400-e29b-41d4-a716-446655440000",
		Email:  req.Email,
		Role:   "learner",
		Token:  "placeholder-token",
	})
}

func handleAuthMe(w http.ResponseWriter, r *http.Request, db *database.Connection, tm *auth.TokenManager) {
	if r.Method != http.MethodGet {
		httputil.BadRequest(w, "Method not allowed")
		return
	}

	// For Phase One, placeholder
	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": "550e8400-e29b-41d4-a716-446655440000",
		"email":   "user@example.com",
		"role":    "learner",
	})
}

// Helper function
func parseJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
