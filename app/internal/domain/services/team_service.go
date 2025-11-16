package services

import (
	"context"
	"log/slog"

	"github.com/qurk0/pr-service/internal/domain/models"
	"github.com/qurk0/pr-service/internal/metrics"
)

type TeamRepo interface {
	CheckTeamExists(ctx context.Context, teamName string) (bool, error)
	CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (models.Team, error)
}

type TeamService struct {
	repo TeamRepo
	log  *slog.Logger
}

func newTeamService(repo TeamRepo, log *slog.Logger) *TeamService {
	return &TeamService{repo: repo, log: log}
}

func (s *TeamService) CheckTeamExists(ctx context.Context, teamName string) (bool, error) {
	return s.repo.CheckTeamExists(ctx, teamName)
}

func (s *TeamService) CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (models.Team, error) {
	const op = "team_service.CreateTeamWithMembers"

	exists, err := s.repo.CheckTeamExists(ctx, teamName)
	if err != nil {
		s.log.Error(op, slog.String("error from repo", err.Error()))
		return models.Team{}, err
	}
	if exists {
		s.log.Warn(op, slog.String("fail: team already exists", teamName))
		return models.Team{}, models.ErrTeamExists
	}

	ids := make(map[string]struct{})
	for _, member := range members {
		if _, ok := ids[member.Id]; ok {
			s.log.Error(op, slog.String("error: duplicated id", member.Id))
			return models.Team{}, models.ErrDuplicatedIds
		}
		ids[member.Id] = struct{}{}
	}
	metrics.TeamsCreated.Inc()
	return s.repo.CreateTeamWithMembers(ctx, teamName, members)
}
