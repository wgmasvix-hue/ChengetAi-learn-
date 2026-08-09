package http

import (
	"encoding/json"
	stdhttp "net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	goredis "github.com/redis/go-redis/v9"
	appauth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/auth"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/config"
	apphealth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/health"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/learners"
	appmw "github.com/wgmasvix-hue/ChengetAi-learn-/internal/middleware"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/teachers"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/users"
	"go.uber.org/zap"
)

func NewRouter(
	cfg config.Config,
	logger *zap.Logger,
	redisClient *goredis.Client,
	authHandler *appauth.Handler,
	usersHandler *users.Handler,
	learnersHandler *learners.Handler,
	teachersHandler *teachers.Handler,
	healthHandler *apphealth.Handler,
) stdhttp.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.Recoverer)
	r.Use(appmw.CORS(allowedOrigins()))
	r.Use(appmw.SecurityHeaders(connectSrc()))
	r.Use(appmw.Logging(logger, "chengetai-learn-api"))

	authMiddleware := appmw.Auth(cfg.JWTSecret, redisClient)
	rateLimiter := appmw.NewIPRateLimiter(float64(cfg.RateLimitRPS), cfg.RateLimitBurst)

	r.Get("/health", healthHandler.Health)
	r.Get("/ready", healthHandler.Ready)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/version", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
			writeJSON(w, stdhttp.StatusOK, map[string]string{
				"name":        "chengetai-learn-api",
				"version":     "0.1.0",
				"environment": cfg.AppEnv,
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.With(rateLimiter.Middleware).Post("/register", authHandler.Register)
			r.With(rateLimiter.Middleware).Post("/login", authHandler.Login)
			r.With(authMiddleware).Post("/logout", authHandler.Logout)
			r.With(authMiddleware).Get("/me", authHandler.Me)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/users/me", usersHandler.Me)
			r.Get("/learners/me", learnersHandler.Me)
			r.Get("/teachers/me", teachersHandler.Me)
		})
	})

	return r
}

func allowedOrigins() []string {
	value := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if value == "" {
		return []string{"*"}
	}
	return strings.Split(value, ",")
}

func connectSrc() string {
	value := strings.TrimSpace(os.Getenv("CSP_CONNECT_SRC"))
	if value == "" {
		return "'self' http://localhost:8080"
	}
	return value
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
