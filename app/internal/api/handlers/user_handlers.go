package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
)

type UserHandler struct {
	userServ UserService
	prServ   PRService
}

func NewUserHandler(uServ UserService, prServ PRService) *UserHandler {
	return &UserHandler{userServ: uServ, prServ: prServ}
}

func (h *UserHandler) SetIsActive(c *fiber.Ctx) error {
	var req dto.SetIsActiveRequest
	_ = c.BodyParser(&req)
	// По документации написано, что 400 не возвращаем - считаю все запросы валидными.
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
	// }

	user, err := h.userServ.SetIsActive(c.UserContext(), req.UserId, req.IsActive)
	if err != nil {
		return writeError(c, err)
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

func (h *UserHandler) GetReview(c *fiber.Ctx) error {
	userID := c.Query("user_id", "")
	// По документации написано, что 400 не возвращаем - считаю все запросы валидными.

	prShortList, err := h.prServ.GetByReviewer(c.UserContext(), userID)
	if err != nil {
		return writeError(c, err)
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
