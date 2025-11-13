package services

type UserRepo interface {
}

type UserService struct {
	repo UserRepo
}

func newUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}
