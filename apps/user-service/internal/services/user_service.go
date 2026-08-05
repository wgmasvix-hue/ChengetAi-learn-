package services

import (
	"context"
	"fmt"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/user-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/user-service/internal/repositories"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// UserService handles user operations
type UserService struct {
	profileRepo *repositories.ProfileRepository
}

// NewUserService creates a new user service
func NewUserService(profileRepo *repositories.ProfileRepository) *UserService {
	return &UserService{
		profileRepo: profileRepo,
	}
}

// GetProfile retrieves a user profile
func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.ProfileResponse, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.ProfileResponse{
		UserID:      profile.UserID,
		FullName:    profile.FullName,
		AvatarURL:   profile.AvatarURL,
		Bio:         profile.Bio,
		PhoneNumber: profile.PhoneNumber,
		Country:     profile.Country,
		Language:    profile.Language,
		Timezone:    profile.Timezone,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}, nil
}

// UpdateProfile updates a user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateProfileRequest) (*domain.ProfileResponse, error) {
	profile := &domain.Profile{
		UserID:      userID,
		FullName:    req.FullName,
		Bio:         &req.Bio,
		AvatarURL:   &req.AvatarURL,
		PhoneNumber: &req.PhoneNumber,
		Timezone:    req.Timezone,
	}

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, fmt.Errorf("updating profile: %w", err)
	}

	// Fetch updated profile
	return s.GetProfile(ctx, userID)
}

// DeleteProfile deletes a user profile
func (s *UserService) DeleteProfile(ctx context.Context, userID string) error {
	return s.profileRepo.Delete(ctx, userID)
}
