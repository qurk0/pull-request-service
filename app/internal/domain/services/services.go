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

func NewServices(db DB) *Services {
	return &Services{
		User: newUserService(db),
		Team: newTeamService(db),
		PR:   newPullRequestService(db),
	}
}
