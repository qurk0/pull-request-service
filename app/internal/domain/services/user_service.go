package services

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type UserRepo interface {
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUserIsActive(ctx context.Context, user *models.User) error
}

type UserService struct {
	repo UserRepo
}

func newUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) SetIsActive(ctx context.Context, userID string, active bool) (*models.User, error) {
	user, err := us.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.IsActive = active

	if err := us.repo.UpdateUserIsActive(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
