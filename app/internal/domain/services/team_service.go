package services

type TeamRepo interface {
}

type TeamService struct {
	repo TeamRepo
}

func newTeamService(repo TeamRepo) *TeamService {
	return &TeamService{repo: repo}
}
