package services

import (
	"context"
)

type TeamRepo interface {
	CheckTeamExists(ctx context.Context, teamName string) error
}

type TeamService struct {
	repo TeamRepo
}

func newTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CheckTeamExists(ctx context.Context, teamName string) error {
	return s.repo.CheckTeamExists(ctx, teamName)
}
