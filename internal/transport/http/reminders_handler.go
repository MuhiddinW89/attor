package transport

import (
	"errors"

	reminderapp "github.com/MuhiddinW89/attor/internal/application/reminders"
	"github.com/MuhiddinW89/attor/internal/reminders"
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

func (h *ReminderHandler) MarkSent(
	c *fiber.Ctx,
) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid reminder id",
			},
		)
	}

	err = h.markSentUC.Execute(
		c.Context(),
		id,
	)

	if err != nil {

		if errors.Is(err, reminders.ErrReminderNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ReminderHandler) List(
	c *fiber.Ctx,
) error {

	reminders, err := h.listUC.Execute(
		c.Context(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(reminders)
}
