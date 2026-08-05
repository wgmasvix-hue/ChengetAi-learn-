package domain

import "time"

// School represents an educational institution
type School struct {
	ID                    string
	Name                  string
	Description           *string
	Country               *string
	Region                *string
	City                  *string
	Website               *string
	Email                 *string
	PhoneNumber           *string
	LogoURL               *string
	Status                string // active, suspended, inactive
	SubscriptionTier      string // free, basic, professional, enterprise
	SubscriptionExpiresAt *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// CreateSchoolRequest represents a school creation request
type CreateSchoolRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Country     string `json:"country" validate:"iso3166_1_alpha2"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Website     string `json:"website"`
	Email       string `json:"email" validate:"email"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateSchoolRequest represents a school update request
type UpdateSchoolRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	LogoURL     string `json:"logo_url"`
}

// SchoolResponse represents a school in API responses
type SchoolResponse struct {
	ID                    string     `json:"id"`
	Name                  string     `json:"name"`
	Description           *string    `json:"description,omitempty"`
	Country               *string    `json:"country,omitempty"`
	Region                *string    `json:"region,omitempty"`
	City                  *string    `json:"city,omitempty"`
	Website               *string    `json:"website,omitempty"`
	Email                 *string    `json:"email,omitempty"`
	PhoneNumber           *string    `json:"phone_number,omitempty"`
	LogoURL               *string    `json:"logo_url,omitempty"`
	Status                string     `json:"status"`
	SubscriptionTier      string     `json:"subscription_tier"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// Curriculum represents a school curriculum
type Curriculum struct {
	ID              string
	SchoolID        string
	Name            string
	Description     *string
	CurriculumType  string // zimsec, cambridge, ib, etc
	CreatedAt       time.Time
}
