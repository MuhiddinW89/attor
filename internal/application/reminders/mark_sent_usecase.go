package reminders

import (
	"context"

	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/google/uuid"
)

type MarkReminderSentUseCase struct {
	reminderService reminders.Service
}

func NewMarkReminderSentUseCase(
	reminderService reminders.Service,
) *MarkReminderSentUseCase {
	return &MarkReminderSentUseCase{
		reminderService: reminderService,
	}
}

func (uc *MarkReminderSentUseCase) Execute(
	ctx context.Context,
	id uuid.UUID,
) error {

	return uc.reminderService.MarkSent(
		ctx,
		id,
	)
}
