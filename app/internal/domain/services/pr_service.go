package services

import (
	"context"
	"math/rand"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type PullRequestRepo interface {
	GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error)
	CreatePR(ctx context.Context, prID, prNamme, authorID string, reviewers []string) (models.PR, error)
	GetPRByID(ctx context.Context, prID string) (models.PR, error)
	ReassignPRReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (models.PR, []string, error)
	GetPRReviewers(ctx context.Context, prID string) ([]string, error)
	MergePR(ctx context.Context, prID string) (models.PR, error)
}

type PullRequestService struct {
	repo  PullRequestRepo
	uServ *UserService
}

func newPullRequestService(repo PullRequestRepo, uServ *UserService) *PullRequestService {
	return &PullRequestService{
		repo:  repo,
		uServ: uServ,
	}
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

func (s *PullRequestService) ReassignPR(ctx context.Context, prID, oldReviewerID string) (models.PR, string, error) {
	pr, err := s.repo.GetPRByID(ctx, prID)
	if err != nil {
		return models.PR{}, "", err
	}

	candidates, err := s.uServ.GetAnotherReviewers(ctx, pr.PRID, oldReviewerID, pr.AuthorID)

	if err != nil {
		return models.PR{}, "", err
	}

	// Тут вылетает ошибка что кандидата нет
	if len(candidates) == 0 {
		return models.PR{}, "", models.ErrNoCandidate
	}

	reviewers := getReviewers(candidates)
	newReviewerID := reviewers[0]

	// Тут транзакция, выплёвывает 1 из 3 возможных ошибок, либо пятисотим
	newPr, newReviewers, err := s.repo.ReassignPRReviewer(ctx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return models.PR{}, "", err
	}

	newPr.AssignedReviewers = newReviewers

	return newPr, newReviewerID, nil
}

func (s *PullRequestService) MergePR(ctx context.Context, prID string) (models.PR, error) {
	newPr, err := s.repo.MergePR(ctx, prID)
	if err != nil {
		return models.PR{}, err
	}

	reviewers, err := s.repo.GetPRReviewers(ctx, prID)
	if err != nil {
		return models.PR{}, err
	}

	newPr.AssignedReviewers = reviewers
	return newPr, nil
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
