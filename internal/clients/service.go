package clients

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input CreateClientInput) (*Client, error)
	List(ctx context.Context, search string) ([]*Client, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Client, error)
	GetByPhone(ctx context.Context, phone string) (*Client, error)
}

type CreateClientInput struct {
	FullName  string
	Phone     string
	Instagram *string
}
