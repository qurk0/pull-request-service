package services

import "log/slog"

type DB interface {
	UserRepo
	TeamRepo
	PullRequestRepo
}

type Services struct {
	User *UserService
	Team *TeamService
	PR   *PullRequestService
}

func NewServices(uRepo UserRepo, tRepo TeamRepo, prRepo PullRequestRepo, logger *slog.Logger) *Services {
	return &Services{
		User: newUserService(uRepo, logger),
		Team: newTeamService(tRepo, logger),
		PR:   newPullRequestService(prRepo, newUserService(uRepo, logger), logger),
	}
}
