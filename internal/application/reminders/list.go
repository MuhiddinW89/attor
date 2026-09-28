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
	status string,
) ([]*ReminderListItem, error) {

	items, err := uc.reminderService.ListDetailed(
		ctx,
		status,
	)
	if err != nil {
		return nil, err
	}

	result := make([]*ReminderListItem, 0, len(items))

	for _, item := range items {
		result = append(result, &ReminderListItem{
			ID:          item.ID.String(),
			ClientID:    item.ClientID.String(),
			ClientName:  item.ClientName,
			Phone:       item.Phone,
			PerfumeName: item.PerfumeName,
			VolumeML:    item.VolumeML,
			Price:       item.Price,
			Comment:     &item.Comment,
			SaleDate:    item.SaleDate,
			ReminderAt:  item.ReminderAt,
			IsSent:      item.IsSent,
		})
	}

	return result, nil
}
