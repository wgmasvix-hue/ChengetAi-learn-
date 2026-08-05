package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID               string
	Email            string
	Username         string
	PasswordHash     string
	FullName         string
	AvatarURL        *string
	Bio              *string
	PhoneNumber      *string
	Country          *string
	Language         string
	Timezone         string
	Status           string
	EmailVerified    bool
	EmailVerifiedAt  *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	Country  string `json:"country" validate:"iso3166_1_alpha2"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
