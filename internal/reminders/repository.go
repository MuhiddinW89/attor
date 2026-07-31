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

	MarkCompleted(
		ctx context.Context,
		id uuid.UUID,
	) error

	GetNearestByClientID(
		ctx context.Context,
		clientID uuid.UUID,
	) (*Reminder, error)
}
