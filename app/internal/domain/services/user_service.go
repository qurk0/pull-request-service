package services

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type UserRepo interface {
	GetUser(ctx context.Context, userID string) (models.User, error)
	UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error
	GetTeamMembers(ctx context.Context, teamName string) ([]models.TeamMember, error)
}

type UserService struct {
	repo UserRepo
}

func newUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (models.User, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	if err := s.repo.UpdateUserIsActive(ctx, user.Id, active); err != nil {
		return models.User{}, err
	}

	user.IsActive = active
	return user, nil
}

func (s *UserService) GetTeamMembers(ctx context.Context, teamName string) ([]models.TeamMember, error) {
	return s.repo.GetTeamMembers(ctx, teamName)
}
