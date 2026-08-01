package reminders

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(
	ctx context.Context,
	input CreateReminderInput,
) (*Reminder, error) {

	now := time.Now()

	reminder := &Reminder{
		ID:         uuid.New(),
		ClientID:   input.ClientID,
		SaleID:     input.SaleID,
		ReminderAt: input.ReminderAt,
		IsSent:     false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repository.Create(ctx, reminder); err != nil {
		return nil, err
	}

	return reminder, nil
}

func (s *service) ListPending(
	ctx context.Context,
) ([]*Reminder, error) {
	return s.repository.ListPending(ctx)
}

func (s *service) GetNearestByClientID(
	ctx context.Context,
	clientID uuid.UUID,
) (*Reminder, error) {

	return s.repository.GetNearestByClientID(
		ctx,
		clientID,
	)
}

func (s *service) MarkSent(
	ctx context.Context,
	id uuid.UUID,
) error {

	if id == uuid.Nil {
		return ErrInvalidReminderID
	}

	return s.repository.MarkSent(
		ctx,
		id,
	)
}

func (s *service) ListPendingDetailed(
	ctx context.Context,
) ([]*ReminderListItem, error) {

	return s.repository.ListPendingDetailed(ctx)
}
