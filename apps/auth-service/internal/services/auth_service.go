package services

import (
	"context"
	"fmt"
	"time"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/auth-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/auth-service/internal/repositories"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/crypto"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/jwt"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.SessionRepository
	jwtManager  *jwt.Manager
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repositories.UserRepository,
	sessionRepo *repositories.SessionRepository,
	jwtManager *jwt.Manager,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtManager:  jwtManager,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.LoginResponse, error) {
	// Hash password
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternal("hashing password failed", err)
	}

	// Create user
	user := &domain.User{
		Email:       req.Email,
		Username:    req.Username,
		PasswordHash: passwordHash,
		FullName:    req.FullName,
		Country:     &req.Country,
		Language:    "en",
		Timezone:    "UTC",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, errors.NewInternal("generating tokens failed", err)
	}

	// Create session
	session := &repositories.Session{
		UserID:           user.ID,
		TokenHash:        crypto.HashToken(tokenPair.AccessToken),
		RefreshTokenHash: crypto.HashToken(tokenPair.RefreshToken),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		RefreshExpiresAt: timePtr(time.Now().Add(7 * 24 * time.Hour)),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, errors.NewInternal("creating session failed", err)
	}

	return &domain.LoginResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FullName:     user.FullName,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()),
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewInvalidPassword()
	}

	// Verify password
	if err := crypto.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.NewInvalidPassword()
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, errors.NewInternal("generating tokens failed", err)
	}

	// Create session
	session := &repositories.Session{
		UserID:           user.ID,
		TokenHash:        crypto.HashToken(tokenPair.AccessToken),
		RefreshTokenHash: crypto.HashToken(tokenPair.RefreshToken),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		RefreshExpiresAt: timePtr(time.Now().Add(7 * 24 * time.Hour)),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, errors.NewInternal("creating session failed", err)
	}

	return &domain.LoginResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FullName:     user.FullName,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()),
	}, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
	// Verify refresh token
	claims, err := s.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.NewInvalidToken()
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.NewUnauthorized("user not found", err)
	}

	// Generate new token pair
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, errors.NewInternal("generating tokens failed", err)
	}

	// Create new session
	session := &repositories.Session{
		UserID:           user.ID,
		TokenHash:        crypto.HashToken(tokenPair.AccessToken),
		RefreshTokenHash: crypto.HashToken(tokenPair.RefreshToken),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		RefreshExpiresAt: timePtr(time.Now().Add(7 * 24 * time.Hour)),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, errors.NewInternal("creating session failed", err)
	}

	return &domain.LoginResponse{
		UserID:       user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FullName:     user.FullName,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()),
	}, nil
}

// Logout logs out a user
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Revoke(ctx, sessionID)
}

// VerifyToken verifies an access token
func (s *AuthService) VerifyToken(ctx context.Context, token string) (*jwt.Claims, error) {
	// Verify token signature
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return nil, errors.NewInvalidToken()
	}

	// Verify session exists
	tokenHash := crypto.HashToken(token)
	_, err = s.sessionRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}
