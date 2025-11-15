package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
	"github.com/qurk0/pr-service/internal/domain/models"
	"github.com/qurk0/pr-service/internal/domain/services"
)

type Router struct {
	User *UserHandler
	Team *TeamHandler
	PR   *PRHandler
}

func NewRouter(servs *services.Services) *Router {
	return &Router{
		User: NewUserHandler(servs.User, servs.PR),
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

func writeError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, models.ErrTimeout):
		return c.Status(fiber.StatusServiceUnavailable).SendString("timeout reached")

	case errors.Is(err, models.ErrCanceled):
		return c.Status(fiber.StatusServiceUnavailable).SendString("request canceled")

	case errors.Is(err, models.ErrInternal):
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")

	case errors.Is(err, models.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: dto.HttpError{
				Code:    dto.ErrCodeNotFound,
				Message: "not found",
			},
		})

	default:
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}
}
