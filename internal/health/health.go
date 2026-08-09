package health

import (
	"context"
	"net/http"
	"time"

	"github.com/wgmasvix-hue/chengetai-learn/internal/database"
	"github.com/wgmasvix-hue/chengetai-learn/internal/http"
	"github.com/wgmasvix-hue/chengetai-learn/internal/redis"
)

// Handler handles health check requests
type Handler struct {
	db    *database.Connection
	redis *redis.Client
}

// NewHandler creates a new health handler
func NewHandler(db *database.Connection, redisClient *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redisClient,
	}
}

// Health handles the /health endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := http.HealthResponse{
		Status: "ok",
	}
	http.WriteJSON(w, http.StatusOK, resp)
}

// Ready handles the /ready endpoint (checks all dependencies)
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	services := make(map[string]http.ServiceStatus)

	// Check database
	dbStatus := http.ServiceStatus{Status: "ok"}
	if err := h.db.Health(ctx); err != nil {
		dbStatus.Status = "error"
		dbStatus.Error = err.Error()
	}
	services["database"] = dbStatus

	// Check Redis
	redisStatus := http.ServiceStatus{Status: "ok"}
	if err := h.redis.Health(ctx); err != nil {
		redisStatus.Status = "error"
		redisStatus.Error = err.Error()
	}
	services["redis"] = redisStatus

	// Determine overall status
	overallStatus := "ready"
	statusCode := http.StatusOK

	for _, status := range services {
		if status.Status == "error" {
			overallStatus = "not_ready"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	resp := http.ReadinessResponse{
		Status:   overallStatus,
		Services: services,
	}

	http.WriteJSON(w, statusCode, resp)
}

// Version handles the /version endpoint
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	resp := http.VersionResponse{
		Name:        "chengetai-learn-api",
		Version:     "0.1.0",
		Environment: "development",
		GoVersion:   "1.26",
		BuildTime:   time.Now().Format(time.RFC3339),
	}

	http.WriteJSON(w, http.StatusOK, resp)
}
