package services

import (
	"context"
	"math/rand"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type PullRequestRepo interface {
	GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error)
	CreatePR(ctx context.Context, prID, prNamme, authorID string, reviewers []string) (models.PR, error)
}

type PullRequestService struct {
	repo  PullRequestRepo
	uServ *UserService
}

func newPullRequestService(repo PullRequestRepo) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (s *PullRequestService) GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error) {
	return s.repo.GetByReviewer(ctx, userID)
}

func (s *PullRequestService) CreatePR(ctx context.Context, prID, prNamme, authorID string) (models.PR, error) {
	user, err := s.uServ.GetUser(ctx, authorID)
	if err != nil {
		return models.PR{}, err
	}

	// Юзер есть - берём до 2х ревьюверов
	candidates, err := s.uServ.GetReviewers(ctx, user.Id, user.TeamName)
	if err != nil {
		return models.PR{}, err
	}

	reviewers := getReviewers(candidates)

	return s.repo.CreatePR(ctx, prID, prNamme, authorID, reviewers)
}

func getReviewers(candidates []string) []string {
	if len(candidates) <= 2 {
		return candidates
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:2]
}
