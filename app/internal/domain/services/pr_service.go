package services

type PullRequestRepo interface {
}

type PullRequestService struct {
	repo PullRequestRepo
}

func newPullRequestService(repo PullRequestRepo) *PullRequestService {
	return &PullRequestService{repo: repo}
}
