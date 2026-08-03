package reminders

import (
	"context"

	domain "github.com/MuhiddinW89/attor/internal/reminders"
)

type ListUseCase struct {
	reminderService domain.Service
}

func NewListUseCase(
	reminderService domain.Service,
) *ListUseCase {
	return &ListUseCase{
		reminderService: reminderService,
	}
}

func (uc *ListUseCase) Execute(
	ctx context.Context,
) ([]*domain.ReminderListItem, error) {

	return uc.reminderService.ListPendingDetailed(ctx)
}