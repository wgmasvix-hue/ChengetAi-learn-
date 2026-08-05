package domain

import "time"

// Profile represents a user profile
type Profile struct {
	UserID    string
	FullName  string
	AvatarURL *string
	Bio       *string
	PhoneNumber *string
	Country   *string
	Language  string
	Timezone  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpdateProfileRequest represents a profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number"`
	Timezone string `json:"timezone"`
}

// ProfileResponse represents a profile response
type ProfileResponse struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Country   *string `json:"country,omitempty"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
