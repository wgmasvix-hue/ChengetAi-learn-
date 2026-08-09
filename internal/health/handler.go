package health

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"time"
)

type Handler struct {
	dbPing    func(context.Context) error
	redisPing func(context.Context) error
}

func NewHandler(dbPing, redisPing func(context.Context) error) *Handler {
	return &Handler{dbPing: dbPing, redisPing: redisPing}
}

func (h *Handler) Health(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
	writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Ready(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	response := map[string]any{
		"status": "ok",
		"services": map[string]map[string]string{
			"postgres": {"status": "ok"},
			"redis":    {"status": "ok"},
		},
	}
	statusCode := stdhttp.StatusOK
	services := response["services"].(map[string]map[string]string)

	if h.dbPing != nil {
		if err := h.dbPing(ctx); err != nil {
			response["status"] = "degraded"
			services["postgres"] = map[string]string{"status": "down", "message": err.Error()}
			statusCode = stdhttp.StatusServiceUnavailable
		}
	}
	if h.redisPing != nil {
		if err := h.redisPing(ctx); err != nil {
			response["status"] = "degraded"
			services["redis"] = map[string]string{"status": "down", "message": err.Error()}
			statusCode = stdhttp.StatusServiceUnavailable
		}
	}

	writeJSON(w, statusCode, response)
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
