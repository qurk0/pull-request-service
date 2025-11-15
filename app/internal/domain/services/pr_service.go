package services

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type PullRequestRepo interface {
	GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error)
}

type PullRequestService struct {
	repo PullRequestRepo
}

func newPullRequestService(repo PullRequestRepo) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (prs *PullRequestService) GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error) {
	return prs.repo.GetByReviewer(ctx, userID)
}
