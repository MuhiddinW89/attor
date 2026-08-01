package application

import "context"

type RemindersQueryRepository interface {
	ListPending(
		ctx context.Context,
	) ([]ReminderListItem, error)
}
