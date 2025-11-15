package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
	"github.com/qurk0/pr-service/internal/domain/models"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, active bool) (*models.User, error)
	GetReview(ctx context.Context, userID string) ([]*models.PRShort, error)
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

func (uh *UserHandler) SetIsActive(c *fiber.Ctx) error {
	var req dto.SetIsActiveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
	}

	user, err := uh.serv.SetIsActive(c.Context(), req.UserId, req.IsActive)
	if err != nil {
		return c.Status(500).SendString("implement this case (uh.SetIsActive)")
	}

	respUser := dto.SetIsActiveResponse{
		UserId:   user.Id,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		User dto.SetIsActiveResponse `json:"user"`
	}{
		User: respUser,
	})
}

func (uh *UserHandler) GetReview(c *fiber.Ctx) error {
	userID := c.Query("user_id", "")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("invalid query params")
	}

	prShortList, err := uh.serv.GetReview(c.Context(), userID)
	if err != nil {
		return c.Status(500).SendString("implement this case (uh.GetReview)")
	}

	reqPRShortList := make([]dto.PRShort, 0, len(prShortList))
	for _, pr := range prShortList {
		reqPRShortList = append(reqPRShortList, dto.PRShort{
			ID:       pr.PRID,
			Name:     pr.PRName,
			AuthorID: pr.AuthorID,
			Status:   string(pr.Status),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.GetReviewResponse{
		UserID:        userID,
		RequestsShort: reqPRShortList,
	})
}
