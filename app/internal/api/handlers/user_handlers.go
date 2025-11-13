package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
)

type UserService interface {
	SetIsActive(ctx context.Context, username string, active bool) error
	GetReview(ctx context.Context, username string) error
}

type UserHandler struct {
	serv UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{serv: service}
}

/*
	app.Get("/users/getReview")
	app.Post("/users/setIsActive")
*/

// func (uh *UserHandler) GetReview(c *fiber.Ctx) error {

// }

func (uh *UserHandler) SetIsActive(c *fiber.Ctx) error {
	var req dto.SetIsActiveRequest
	if err := c.BodyParser(&req); err != nil {
		return dto.ReturnError("bad_request", "invalid json body")
	}
}
