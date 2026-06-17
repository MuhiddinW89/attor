package clients

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input CreateClientInput) (*Client, error)
	List(ctx context.Context) ([]*Client, error)
	GetByID(ctx context.Context, id uuid.UUID ) (*Client, error)
}

type CreateClientInput struct {
	FullName  string
	Phone     string
	Instagram *string
}
