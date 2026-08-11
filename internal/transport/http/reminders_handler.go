package transport

import (
	reminderapp "github.com/MuhiddinW89/attor/internal/application/reminders"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReminderHandler struct {
	markSentUC *reminderapp.MarkReminderSentUseCase
	listUC     *reminderapp.ListUseCase
}

func NewReminderHandler(
	markSentUC *reminderapp.MarkReminderSentUseCase,
	listUC *reminderapp.ListUseCase,
) *ReminderHandler {
	return &ReminderHandler{
		markSentUC: markSentUC,
		listUC:     listUC,
	}
}

func (h *ReminderHandler) List(c *fiber.Ctx) error {
	status := c.Query("status", "all")

	switch status {
	case "all", "urgent", "overdue", "completed":
		// valid status

	default:
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid reminder status",
			},
		)
	}

	result, err := h.listUC.Execute(
		c.Context(),
		status,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(result)
}

func (h *ReminderHandler) MarkSent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid reminder id",
			},
		)
	}

	if err := h.markSentUC.Execute(
		c.Context(),
		id,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
