package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid session")
)

const SessionLifetime = 12 * time.Hour

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*User, string, error)
	CreateSession(ctx context.Context, tokenHash, csrfHash []byte, userID uuid.UUID, expiresAt time.Time) error
	GetSession(ctx context.Context, tokenHash []byte) (*Session, error)
	DeleteSession(ctx context.Context, tokenHash []byte) error
}

type LoginResult struct {
	User      *User
	Token     string
	CSRFToken string
	ExpiresAt time.Time
}

type AuthService struct {
	repository AuthRepository
	now        func() time.Time
}

func NewAuthService(repository AuthRepository) *AuthService {
	return &AuthService{repository: repository, now: time.Now}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, passwordHash, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrInvalidCredentials) {
		return nil, err
	}
	if user == nil {
		passwordHash = string(invalidUserPasswordHash)
	}
	passwordMatches := VerifyPassword(passwordHash, password)
	if !passwordMatches || user == nil {
		return nil, ErrInvalidCredentials
	}

	token, tokenHash, err := GenerateOpaqueToken()
	if err != nil {
		return nil, err
	}
	csrfToken, csrfHash, err := GenerateOpaqueToken()
	if err != nil {
		return nil, err
	}
	expiresAt := s.now().UTC().Add(SessionLifetime)
	if err := s.repository.CreateSession(ctx, tokenHash, csrfHash, user.ID, expiresAt); err != nil {
		return nil, err
	}
	return &LoginResult{User: user, Token: token, CSRFToken: csrfToken, ExpiresAt: expiresAt}, nil
}

func (s *AuthService) Authenticate(ctx context.Context, token string) (*Session, error) {
	if len(token) < 32 || len(token) > 128 {
		return nil, ErrInvalidSession
	}
	session, err := s.repository.GetSession(ctx, HashToken(token))
	if err != nil {
		return nil, err
	}
	if session == nil || !session.ExpiresAt.After(s.now()) {
		return nil, ErrInvalidSession
	}
	return session, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	if len(token) < 32 || len(token) > 128 {
		return ErrInvalidSession
	}
	return s.repository.DeleteSession(ctx, HashToken(token))
}
