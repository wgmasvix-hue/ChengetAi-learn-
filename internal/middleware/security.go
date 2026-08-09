package middleware

import (
	"encoding/json"
	"math"
	stdhttp "net/http"
	"strings"
	"sync"
	"time"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

type IPRateLimiter struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	bucket map[string]*bucket
}

func SecurityHeaders(connectSrc string) func(stdhttp.Handler) stdhttp.Handler {
	if strings.TrimSpace(connectSrc) == "" {
		connectSrc = "'self'"
	}
	policy := "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; connect-src " + connectSrc
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Content-Security-Policy", policy)
			next.ServeHTTP(w, r)
		})
	}
}

func CORS(origins []string) func(stdhttp.Handler) stdhttp.Handler {
	allowed := make(map[string]struct{}, len(origins))
	allowAll := false
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "*" {
			allowAll = true
		}
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
				}
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			}
			if r.Method == stdhttp.MethodOptions {
				w.WriteHeader(stdhttp.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NewIPRateLimiter(ratePerSecond float64, burst int) *IPRateLimiter {
	if ratePerSecond <= 0 {
		ratePerSecond = 1.0 / 6.0
	}
	if burst <= 0 {
		burst = 10
	}
	return &IPRateLimiter{
		rate:   ratePerSecond,
		burst:  float64(burst),
		bucket: make(map[string]*bucket),
	}
}

func (l *IPRateLimiter) Middleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		ip := clientIP(r)
		if !l.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(stdhttp.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]string{
					"code":    "RATE_LIMITED",
					"message": "Too many requests",
				},
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (l *IPRateLimiter) allow(ip string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.bucket[ip]
	if !ok {
		l.bucket[ip] = &bucket{tokens: l.burst - 1, lastRefill: now}
		return true
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = math.Min(l.burst, b.tokens+elapsed*l.rate)
	b.lastRefill = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func clientIP(r *stdhttp.Request) string {
	if forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			return strings.TrimSpace(parts[0])
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	remoteAddr := r.RemoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx > 0 {
		return remoteAddr[:idx]
	}
	return remoteAddr
}
