package sales

import (
	"context"
	"github.com/google/uuid"
)

type service struct {
	repository Repository
}

// Create implements [Service].
func (s *service) Create(ctx context.Context, input CreateSaleInput) (*Sale, error) {
	panic("unimplemented")
}



func NewService(
	repository Repository,
) Service {
	return &service{
		repository: repository,
	}
}
