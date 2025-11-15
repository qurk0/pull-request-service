package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
	"github.com/qurk0/pr-service/internal/domain/models"
)

type TeamHandler struct {
	teamServ TeamService
	userServ UserService
}

func NewTeamHandler(tService TeamService, uService UserService) *TeamHandler {
	return &TeamHandler{
		teamServ: tService,
		userServ: uService,
	}
}

func (h *TeamHandler) GetTeam(c *fiber.Ctx) error {
	teamName := c.Query("team_name", "")
	// По документации написано, что 400 не возвращаем - считаю все запросы валидными.

	exists, err := h.teamServ.CheckTeamExists(c.UserContext(), teamName)
	if err != nil {
		return writeError(c, err)
	}

	if !exists {
		return writeError(c, models.ErrNotFound)
	}

	// Сюда идём когда команда существует, в случае чего мы получим пустой список.
	teamMembers, err := h.userServ.GetTeamMembers(c.UserContext(), teamName)
	if err != nil {
		return writeError(c, err)
	}

	respMembers := make([]dto.Member, 0, len(teamMembers))
	for _, teamMember := range teamMembers {
		respMembers = append(respMembers, dto.Member{
			Id:       teamMember.Id,
			Username: teamMember.Username,
			IsActive: teamMember.IsActive,
		})
	}

	resp := dto.TeamResponse{
		TeamName: teamName,
		Members:  respMembers,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *TeamHandler) AddTeam(c *fiber.Ctx) error {
	var req dto.AddTeamRequest
	_ = c.BodyParser(&req)
	//	По документации написано, что 400 не возвращаем - считаю все запросы валидными.
	//  if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
	// }

	members := membersDTOtoModel(req.Members)

	team, err := h.teamServ.CreateTeamWithMembers(c.UserContext(), req.TeamName, members)
	if err != nil {
		return writeError(c, err)
	}

	respMembers := membersModelToDTO(team.TeamMembers)

	return c.Status(fiber.StatusCreated).JSON(struct {
		Team dto.TeamResponse `json:"team"`
	}{
		Team: dto.TeamResponse{
			TeamName: team.TeamName,
			Members:  respMembers,
		},
	})
}

func membersDTOtoModel(in []dto.Member) []models.TeamMember {
	out := make([]models.TeamMember, 0, len(in))
	for _, member := range in {
		out = append(out, models.TeamMember{
			Id:       member.Id,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return out
}

func membersModelToDTO(in []models.TeamMember) []dto.Member {
	out := make([]dto.Member, 0, len(in))
	for _, member := range in {
		out = append(out, dto.Member{
			Id:       member.Id,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return out
}
