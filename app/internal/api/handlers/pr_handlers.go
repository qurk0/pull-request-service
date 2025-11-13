package handlers

import "github.com/qurk0/pr-service/internal/domain/services"

type PRHandler struct {
	serv *services.PullRequestService
}

func NewPRHandler(service *services.PullRequestService) *PRHandler {
	return &PRHandler{serv: service}
}
