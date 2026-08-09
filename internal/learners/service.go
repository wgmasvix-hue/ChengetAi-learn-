package learners

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (*Profile, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	return s.repo.FindByUserID(ctx, id)
}
