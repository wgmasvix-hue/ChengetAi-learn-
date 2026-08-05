package services

import (
	"context"
	"fmt"

	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/school-service/internal/domain"
	"github.com/wgmasvix-hue/ChengetAi-learn-/apps/school-service/internal/repositories"
	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/errors"
)

// SchoolService handles school operations
type SchoolService struct {
	schoolRepo *repositories.SchoolRepository
}

// NewSchoolService creates a new school service
func NewSchoolService(schoolRepo *repositories.SchoolRepository) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
	}
}

// CreateSchool creates a new school
func (s *SchoolService) CreateSchool(ctx context.Context, req *domain.CreateSchoolRequest) (*domain.SchoolResponse, error) {
	school := &domain.School{
		Name:        req.Name,
		Description: &req.Description,
		Country:     &req.Country,
		Region:      &req.Region,
		City:        &req.City,
		Website:     &req.Website,
		Email:       &req.Email,
		PhoneNumber: &req.PhoneNumber,
		Status:      "active",
	}

	if err := s.schoolRepo.Create(ctx, school); err != nil {
		return nil, err
	}

	return s.schoolToResponse(school), nil
}

// GetSchool retrieves a school by ID
func (s *SchoolService) GetSchool(ctx context.Context, schoolID string) (*domain.SchoolResponse, error) {
	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	return s.schoolToResponse(school), nil
}

// UpdateSchool updates a school
func (s *SchoolService) UpdateSchool(ctx context.Context, schoolID string, req *domain.UpdateSchoolRequest) (*domain.SchoolResponse, error) {
	// Get existing school
	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		school.Name = req.Name
	}
	if req.Description != "" {
		school.Description = &req.Description
	}
	if req.Website != "" {
		school.Website = &req.Website
	}
	if req.Email != "" {
		school.Email = &req.Email
	}
	if req.PhoneNumber != "" {
		school.PhoneNumber = &req.PhoneNumber
	}
	if req.LogoURL != "" {
		school.LogoURL = &req.LogoURL
	}

	if err := s.schoolRepo.Update(ctx, school); err != nil {
		return nil, fmt.Errorf("updating school: %w", err)
	}

	return s.schoolToResponse(school), nil
}

// ListSchoolsByCountry lists schools by country
func (s *SchoolService) ListSchoolsByCountry(ctx context.Context, country string, limit, offset int) ([]domain.SchoolResponse, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	schools, err := s.schoolRepo.ListByCountry(ctx, country, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]domain.SchoolResponse, len(schools))
	for i, school := range schools {
		response[i] = *s.schoolToResponse(&school)
	}

	return response, nil
}

// Helper function to convert school to response
func (s *SchoolService) schoolToResponse(school *domain.School) *domain.SchoolResponse {
	return &domain.SchoolResponse{
		ID:                    school.ID,
		Name:                  school.Name,
		Description:           school.Description,
		Country:               school.Country,
		Region:                school.Region,
		City:                  school.City,
		Website:               school.Website,
		Email:                 school.Email,
		PhoneNumber:           school.PhoneNumber,
		LogoURL:               school.LogoURL,
		Status:                school.Status,
		SubscriptionTier:      school.SubscriptionTier,
		SubscriptionExpiresAt: school.SubscriptionExpiresAt,
		CreatedAt:             school.CreatedAt,
		UpdatedAt:             school.UpdatedAt,
	}
}
