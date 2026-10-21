package identity

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

var ErrInitialOwnerExists = errors.New("initial owner already exists")

type Repository interface {
	CreateInitialOwner(ctx context.Context, email, passwordHash, storeName string) (*User, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateInitialOwner(ctx context.Context, email, password, storeName string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || len(email) > 320 {
		return nil, errors.New("invalid owner email")
	}
	if len(password) < 12 || len(password) > 72 {
		return nil, errors.New("owner password must be 12-72 bytes")
	}
	storeName = strings.TrimSpace(storeName)
	if storeName == "" || utf8.RuneCountInString(storeName) > 255 {
		return nil, errors.New("invalid store name")
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return s.repository.CreateInitialOwner(ctx, email, passwordHash, storeName)
}
