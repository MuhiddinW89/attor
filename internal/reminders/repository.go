package reminders

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		reminder *Reminder,
	) error

	ListPending(
		ctx context.Context,
	) ([]*Reminder, error)

	GetNearestByClientID(
		ctx context.Context,
		clientID uuid.UUID,
	) (*Reminder, error)

	MarkSent(
		ctx context.Context,
		id uuid.UUID,
	) error

	ListPendingDetailed(
		ctx context.Context,
	) ([]*ReminderListItem, error)
}
