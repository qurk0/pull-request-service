package handlers

import (
	"github.com/qurk0/pr-service/internal/domain/services"
)

type TeamHandler struct {
	serv *services.TeamService
}

func NewTeamHandler(service *services.TeamService) *TeamHandler {
	return &TeamHandler{serv: service}
}
