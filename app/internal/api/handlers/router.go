package handlers

import (
	"context"
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

type UserService interface {
	SetIsActive(ctx context.Context, userID string, active bool) (models.User, error)
	GetTeamMembers(ctx context.Context, teamName string) ([]models.TeamMember, error)
	GetUser(ctx context.Context, userID string) (models.User, error)
}

type PRService interface {
	GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error)
	CreatePR(ctx context.Context, prID, prNamme, authorID string) (models.PR, error)
	ReassignPR(ctx context.Context, prID, oldReviewerID string) (models.PR, string, error)
}

type TeamService interface {
	CheckTeamExists(ctx context.Context, teamName string) (bool, error)
	CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (models.Team, error)
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
	case errors.Is(err, models.ErrTeamExists):
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.HttpError{
				Code:    dto.ErrCodeTeamExists,
				Message: "team_name already exists",
			},
		})

	case errors.Is(err, models.ErrPRExists):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: dto.HttpError{
				Code:    dto.ErrCodePrExists,
				Message: "PR id already exists",
			},
		})

	case errors.Is(err, models.ErrPRMerged):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: dto.HttpError{
				Code:    dto.ErrCodePrMerged,
				Message: "cannot reassign on merged PR",
			},
		})

	case errors.Is(err, models.ErrNotAssigned):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: dto.HttpError{
				Code:    dto.ErrCodeNotAssigned,
				Message: "reviewer is not assigned to this PR",
			},
		})

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

	case errors.Is(err, models.ErrNoCandidate):
		return c.Status(fiber.StatusConflict).JSON(dto.ErrorResponse{
			Error: dto.HttpError{
				Code:    dto.ErrCodeNoCandidate,
				Message: "no reviewer candidate available",
			},
		})

	default:
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}
}
