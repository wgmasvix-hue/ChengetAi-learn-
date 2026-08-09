package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	appauth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/auth"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/config"
	apphealth "github.com/wgmasvix-hue/ChengetAi-learn-/internal/health"
	apphttp "github.com/wgmasvix-hue/ChengetAi-learn-/internal/http"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/learners"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/teachers"
	"github.com/wgmasvix-hue/ChengetAi-learn-/internal/users"
	"go.uber.org/zap"
)

func TestPublicEndpoints(t *testing.T) {
	cfg := config.Config{AppEnv: "test", JWTSecret: "abcdefghijklmnopqrstuvwxyz123456", RateLimitRPS: 10, RateLimitBurst: 20}
	router := apphttp.NewRouter(
		cfg,
		zap.NewNop(),
		(*goredis.Client)(nil),
		appauth.NewHandler(nil),
		users.NewHandler(nil),
		learners.NewHandler(nil),
		teachers.NewHandler(nil),
		apphealth.NewHandler(func(context.Context) error { return nil }, func(context.Context) error { return nil }),
	)

	cases := []struct {
		name string
		path string
	}{
		{name: "health", path: "/health"},
		{name: "version", path: "/api/v1/version"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d", rec.Code)
			}
		})
	}
}
