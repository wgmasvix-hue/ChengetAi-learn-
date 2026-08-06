package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT token claims with security best practices
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	SchoolID string `json:"school_id,omitempty"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	secret         []byte
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
	issuer         string
	audience       string
}

// NewJWTManager creates a new JWT manager with security configuration
func NewJWTManager(secret string, accessExpiry, refreshExpiry time.Duration, issuer, audience string) (*JWTManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT secret must be at least 32 characters long")
	}
	if accessExpiry == 0 {
		accessExpiry = 15 * time.Minute // Default: 15 minutes
	}
	if refreshExpiry == 0 {
		refreshExpiry = 7 * 24 * time.Hour // Default: 7 days
	}

	return &JWTManager{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
		audience:      audience,
	}, nil
}

// GenerateAccessToken creates a short-lived access token
func (jm *JWTManager) GenerateAccessToken(userID, email, role, schoolID string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		SchoolID: schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    jm.issuer,
			Audience:  jwt.ClaimStrings{jm.audience},
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a long-lived refresh token
func (jm *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(jm.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    jm.issuer,
		Audience:  jwt.ClaimStrings{jm.audience},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jm.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken validates and parses an access token
func (jm *JWTManager) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Additional validation for audience and issuer
	if !claims.VerifyAudience(jm.audience, true) {
		return nil, fmt.Errorf("invalid token audience")
	}
	if !claims.VerifyIssuer(jm.issuer, true) {
		return nil, fmt.Errorf("invalid token issuer")
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses a refresh token
func (jm *JWTManager) ValidateRefreshToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	// Additional validation
	if !claims.VerifyAudience(jm.audience, true) {
		return "", fmt.Errorf("invalid token audience")
	}
	if !claims.VerifyIssuer(jm.issuer, true) {
		return "", fmt.Errorf("invalid token issuer")
	}

	return claims.Subject, nil
}

// RevokeToken should be implemented with a token blacklist (Redis)
// This is a placeholder for the interface
type TokenRevoker interface {
	Revoke(tokenString string, expiry time.Duration) error
	IsRevoked(tokenString string) bool
}
