package services

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type UserRepo interface {
	GetUser(ctx context.Context, userID string) (models.User, error)
	UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error
	GetTeamMembers(ctx context.Context, teamName string) ([]models.TeamMember, error)
	GetReviewers(ctx context.Context, userID, teamName string) ([]string, error)
	GetAnotherReviewers(ctx context.Context, prID, oldReviewerID, authorID string) ([]string, error)
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

func (s *UserService) GetUser(ctx context.Context, userID string) (models.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *UserService) GetReviewers(ctx context.Context, userID, teamName string) ([]string, error) {
	return s.repo.GetReviewers(ctx, userID, teamName)
}

func (s *UserService) GetAnotherReviewers(ctx context.Context, prID, oldReviewerID, authorID string) ([]string, error) {
	return s.repo.GetAnotherReviewers(ctx, prID, oldReviewerID, authorID)
}
