package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	stdhttp "net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	goredis "github.com/redis/go-redis/v9"
)

type authClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func Auth(secret string, client *goredis.Client) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			tokenString, err := TokenFromHeader(r.Header.Get("Authorization"))
			if err != nil {
				writeUnauthorized(w, "Authentication token is required")
				return
			}

			claims := &authClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				writeUnauthorized(w, "Authentication token is invalid")
				return
			}

			if client != nil {
				exists, redisErr := client.Exists(r.Context(), sessionKey(tokenString)).Result()
				if redisErr != nil || exists == 0 {
					writeUnauthorized(w, "Authentication token is invalid")
					return
				}
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, userRoleKey, claims.Role)
			ctx = context.WithValue(ctx, userEmailKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TokenFromHeader(header string) (string, error) {
	if strings.TrimSpace(header) == "" {
		return "", errors.New("missing authorization header")
	}
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return "", errors.New("invalid authorization scheme")
	}
	token := strings.TrimSpace(header[7:])
	if token == "" {
		return "", errors.New("missing token")
	}
	return token, nil
}

func sessionKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return "auth:session:" + hex.EncodeToString(sum[:])
}

func writeUnauthorized(w stdhttp.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stdhttp.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}
