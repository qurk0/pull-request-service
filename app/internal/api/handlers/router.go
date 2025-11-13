package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/domain/services"
)

type Router struct {
	User *UserHandler
	Team *TeamHandler
	PR   *PRHandler
}

func NewRouter(servs *services.Services) *Router {
	return &Router{
		User: NewUserHandler(servs.User),
		Team: NewTeamHandler(servs.Team),
		PR:   NewPRHandler(servs.PR),
	}
}

func (r *Router) RegRoutes(app *fiber.App) {
	app.Post("/team/add")
	app.Get("/team/get")

	app.Get("/users/getReview")
	app.Post("/users/setIsActive")

	app.Post("/pullRequest/create")
	app.Post("/pullRequest/merge")
	app.Post("/pullRequest/reassign")
}
