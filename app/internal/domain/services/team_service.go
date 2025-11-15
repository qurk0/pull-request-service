package services

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

type TeamRepo interface {
	CheckTeamExists(ctx context.Context, teamName string) (bool, error)
	CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (models.Team, error)
}

type TeamService struct {
	repo TeamRepo
}

func newTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CheckTeamExists(ctx context.Context, teamName string) (bool, error) {
	return s.repo.CheckTeamExists(ctx, teamName)
}

func (s *TeamService) CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (models.Team, error) {
	return s.repo.CreateTeamWithMembers(ctx, teamName, members)
}
