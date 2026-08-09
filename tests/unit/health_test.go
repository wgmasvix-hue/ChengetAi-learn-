package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/health"
)

func TestHealthHandler(t *testing.T) {
	h := health.NewHandler(nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
