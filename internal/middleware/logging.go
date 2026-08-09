package middleware

import (
	"context"
	stdhttp "net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
	userRoleKey  contextKey = "user_role"
	userEmailKey contextKey = "user_email"
)

type statusRecorder struct {
	stdhttp.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Logging(logger *zap.Logger, service string) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			r = r.WithContext(ctx)
			w.Header().Set("X-Request-ID", requestID)

			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: stdhttp.StatusOK}
			next.ServeHTTP(recorder, r)

			logger.Info("request completed",
				zap.String("service", service),
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", recorder.status),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey).(string)
	return value
}

func UserIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(userIDKey).(string)
	return value
}

func UserRoleFromContext(ctx context.Context) string {
	value, _ := ctx.Value(userRoleKey).(string)
	return value
}

func UserEmailFromContext(ctx context.Context) string {
	value, _ := ctx.Value(userEmailKey).(string)
	return value
}
