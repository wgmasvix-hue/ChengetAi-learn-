package users

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

func (s *Service) GetByID(ctx context.Context, userID string) (*User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	return s.repo.FindByID(ctx, id)
}
