package clients

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type clientService struct {
	repo Repository
}

func NewService(
	repo Repository,
) Service {
	return &clientService{
		repo: repo,
	}
}

func (s *clientService) Create(
	ctx context.Context,
	input CreateClientInput,
) (*Client, error) {

	existingClient, err := s.repo.GetByPhone(
		ctx,
		input.Phone,
	)

	if err == nil && existingClient != nil {
		return nil, ErrClientAlreadyExists
	}

	if err != nil && !errors.Is(err, ErrClientNotFound) {
		return nil, err
	}

	client := &Client{
		ID:        uuid.New(),
		FullName:  input.FullName,
		Phone:     input.Phone,
		Instagram: input.Instagram,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.Create(
		ctx,
		client,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}


func (s *clientService) List(
	ctx context.Context,
	)([]*Client, error){
	return s.repo.List(ctx)
}

func (s *clientService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Client, error) {

	return s.repo.GetByID(
		ctx,
		id,
	)
}