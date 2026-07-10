package sales

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repository Repository
}

func NewService(
	repository Repository,
) Service {
	return &service{
		repository: repository,
	}
}

// GetByID implements [Service].
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Sale, error) {
	panic("unimplemented")
}

// ListByClientID implements [Service].
func (s *service) ListByClientID(ctx context.Context, clientID uuid.UUID) ([]*Sale, error) {
	panic("unimplemented")
}

// Create implements [Service].
func (s *service) Create(ctx context.Context, input CreateSaleInput) (*Sale, error) {
	if input.Price <= 0 {
		return nil, ErrInvalidSalePrice
	}

	if input.VolumeML <= 0 {
		return nil, ErrInvalidVolume
	}

	if input.PerfumeName == "" {
		return nil, ErrPerfumeNameRequired
	}

	sale := &Sale{
		ID:	uuid.New(),
		ClientID:	input.ClientID,
		PerfumeName: input.PerfumeName,
		VolumeML:	input.VolumeML,
		Price:	input.Price,
		Comment:	input.Comment,
		SaleDate:	input.SaleDate,
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
	}

	if err := s.repository.Create(ctx, sale); err !=nil {
		return nil, err
	}

	return sale, nil
}

