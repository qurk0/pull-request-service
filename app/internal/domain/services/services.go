package services

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

func NewServices(uRepo UserRepo, tRepo TeamRepo, prRepo PullRequestRepo) *Services {
	return &Services{
		User: newUserService(uRepo),
		Team: newTeamService(tRepo),
		PR:   newPullRequestService(prRepo, newUserService(uRepo)),
	}
}
