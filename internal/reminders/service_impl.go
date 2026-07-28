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
		ID:          uuid.New(),
		ClientID:    input.ClientID,
		SaleID:      input.SaleID,
		ReminderAt:  input.ReminderAt,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
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

func (s *service) MarkCompleted(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repository.MarkCompleted(ctx, id)
}
