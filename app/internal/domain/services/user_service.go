package services

import (
	"context"
	"log/slog"

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
	log  *slog.Logger
}

func newUserService(repo UserRepo, log *slog.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) (models.User, error) {
	const op = "user_service.SetIsActive"

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		s.log.Error(op, slog.String("error from repo", err.Error()))
		return models.User{}, err
	}

	s.log.Debug(op, slog.String("got user with id", user.Id))

	if err := s.repo.UpdateUserIsActive(ctx, user.Id, active); err != nil {
		s.log.Error(op, slog.String("error from repo", err.Error()))
		return models.User{}, err
	}

	user.IsActive = active
	s.log.Debug(op, slog.String("success", "user activity changed"))
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
