package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qurk0/pr-service/internal/api/dto"
	"github.com/qurk0/pr-service/internal/domain/models"
)

type PRHandler struct {
	prServ PRService
}

func NewPRHandler(prServ PRService) *PRHandler {
	return &PRHandler{
		prServ: prServ,
	}
}

func (h *PRHandler) CreatePR(c *fiber.Ctx) error {
	var req dto.CreatePRRequest
	_ = c.BodyParser(&req)
	// По документации написано, что 400 не возвращаем - считаю все запросы валидными.
	// if err != nil {
	// 	return writeError(c, err)
	// }

	pr, err := h.prServ.CreatePR(c.UserContext(), req.PRID, req.PRNamme, req.AuthorID)
	if err != nil {
		return writeError(c, err)
	}

	respPr := prModelToDTO(pr)
	return c.Status(fiber.StatusCreated).JSON(struct {
		PullRequest dto.PR `json:"pr"`
	}{
		PullRequest: respPr,
	})
}

func (h *PRHandler) Reassign(c *fiber.Ctx) error {
	var req dto.ReassignPRRequest
	_ = c.BodyParser(&req)
	// По документации написано, что 400 не возвращаем - считаю все запросы валидными.

	pr, newReviewerId, err := h.prServ.ReassignPR(c.UserContext(), req.PRID, req.OldReviewerID)
	if err != nil {
		return writeError(c, err)
	}

	respPr := prModelToDTO(pr)
	return c.Status(fiber.StatusOK).JSON(dto.ReassignPRResponse{
		Pr:         respPr,
		ReplacedBy: newReviewerId,
	})
}

func (h *PRHandler) Merge(c *fiber.Ctx) error {
	var req dto.PRMergeRequest
	_ = c.BodyParser(&req)
	// По документации написано, что 400 не возвращаем - считаю все запросы валидными.

	pr, err := h.prServ.MergePR(c.UserContext(), req.PRID)
	if err != nil {
		return writeError(c, err)
	}

	respPr := prModelToDTO(pr)
	return c.Status(fiber.StatusOK).JSON(struct {
		Pr dto.PR `json:"pr"`
	}{Pr: respPr})
}

func prModelToDTO(in models.PR) dto.PR {
	return dto.PR{
		PRID:              in.PRID,
		PRName:            in.PRName,
		AuthorID:          in.AuthorID,
		Status:            string(in.Status),
		AssignedReviewers: in.AssignedReviewers,
		CreatedAt:         in.CreatedAt,
		MergedAt:          in.MergedAt,
	}
}
