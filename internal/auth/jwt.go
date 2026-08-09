package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Claims represents JWT claims
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
}

// TokenManager handles JWT token operations
type TokenManager struct {
	secret    string
	expiryDuration time.Duration
	issuer    string
}

// NewTokenManager creates a new token manager
func NewTokenManager(secret string, expiryString, issuer string) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("secret must be at least 32 characters long")
	}

	// Parse expiry duration
	duration := 15 * time.Minute
	if expiryString != "" {
		var err error
		duration, err = time.ParseDuration(expiryString)
		if err != nil {
			return nil, fmt.Errorf("invalid token expiry duration: %w", err)
		}
	}

	return &TokenManager{
		secret:         secret,
		expiryDuration: duration,
		issuer:         issuer,
	}, nil
}

// GenerateToken creates a new JWT token
func (tm *TokenManager) GenerateToken(userID, email, role string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(tm.expiryDuration).Unix(),
		Issuer:    tm.issuer,
		Subject:   userID,
	}

	// Create header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	// Encode header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode claims
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signature
	message := headerEncoded + "." + claimsEncoded
	h := hmac.New(sha256.New, []byte(tm.secret))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	token := message + "." + signature
	return token, nil
}

// ValidateToken validates and parses a JWT token
func (tm *TokenManager) ValidateToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerEncoded := parts[0]
	claimsEncoded := parts[1]
	signatureEncoded := parts[2]

	// Verify signature
	message := headerEncoded + "." + claimsEncoded
	h := hmac.New(sha256.New, []byte(tm.secret))
	h.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signatureEncoded), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	// Decode claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(claimsEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token has expired")
	}

	// Verify issuer
	if claims.Issuer != tm.issuer {
		return nil, fmt.Errorf("invalid token issuer")
	}

	return &claims, nil
}
