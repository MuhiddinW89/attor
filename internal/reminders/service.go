package reminders

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(
		ctx context.Context,
		input CreateReminderInput,
	) (*Reminder, error)

	ListPending(
		ctx context.Context,
	) ([]*Reminder, error)

	MarkSent(
		ctx context.Context,
		id uuid.UUID,
	) error

	GetNearestByClientID(
		ctx context.Context,
		clientID uuid.UUID,
	) (*Reminder, error)

	ListPendingDetailed(
		ctx context.Context,
	) ([]*ReminderListItem, error)
}

type CreateReminderInput struct {
	ClientID   uuid.UUID
	SaleID     uuid.UUID
	ReminderAt time.Time
}
