package sales

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	
	Create(
		ctx context.Context,
		input CreateSaleInput,
	) (*Sale, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*Sale, error)
	
	ListByClientID(
		ctx context.Context,
		clientID uuid.UUID,
	) ([]*Sale, error)
}

type CreateSaleInput struct{
	ClientID uuid.UUID
	PerfumeName string
	VolumeML int
	Price float64
	Comment *string
	SaleDate time.Time
}