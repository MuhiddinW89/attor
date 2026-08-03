package reminders

import (
	"context"
)

import domain "github.com/MuhiddinW89/attor/internal/reminders"

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
) ([]*ReminderListItem, error) {

	items, err := uc.reminderService.ListPendingDetailed(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*ReminderListItem, 0, len(items))

	for _, item := range items {

		result = append(result, &ReminderListItem{
			ID:          item.ID,
			ClientID:    item.ClientID,
			ClientName:  item.ClientName,
			Phone:       item.Phone,
			PerfumeName: item.PerfumeName,
			VolumeML:    item.VolumeML,
			Comment:     &item.Comment,
			ReminderAt:  item.ReminderAt,
		})
	}

	return result, nil
}