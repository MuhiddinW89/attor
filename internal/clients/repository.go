package clients

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*Client, error)
	GetByPhone(ctx context.Context, phone string) (*Client, error)
	List(ctx context.Context) ([]*Client, error)
}
