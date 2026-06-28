package sales

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		sale *Sale,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*Sale, error)

	ListByClientID(
		ctx context.Context,
		clientID uuid.UUID,
	) ([]*Sale, error)
}