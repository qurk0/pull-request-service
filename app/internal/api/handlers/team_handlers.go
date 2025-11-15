package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
)

type TeamHandler struct {
	teamServ TeamService
	userServ UserService
}

func NewTeamHandler(tService TeamService) *TeamHandler {
	return &TeamHandler{teamServ: tService}
}

func (h *TeamHandler) GetTeam(c *fiber.Ctx) error {
	teamName := c.Query("team_name", "")
	// Тут должна быть проверка на наличие teamName, но по документации это поле required

	// Если команды нет - ошибка из сервиса будет либо 404, либо 5хх
	if err := h.teamServ.CheckTeamExists(c.UserContext(), teamName); err != nil {
		return writeError(c, err)
	}

	// Сюда идём когда команда существует, в случае чего мы получим пустой список.
	teamMembers, err := h.userServ.GetTeamMembers(c.UserContext(), teamName)
	if err != nil {
		return writeError(c, err)
	}

	respMembers := make([]*dto.Member, 0, len(teamMembers))
	for _, teamMember := range teamMembers {
		respMembers = append(respMembers, &dto.Member{
			Id:       teamMember.Id,
			Username: teamMember.Username,
			IsActive: teamMember.IsActive,
		})
	}

	resp := &dto.GetTeamResponse{
		TeamName: teamName,
		Members:  respMembers,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
