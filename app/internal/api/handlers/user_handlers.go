package handlers

import "github.com/qurk0/pr-service/internal/domain/services"

type UserHandler struct {
	serv *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{serv: service}
}
