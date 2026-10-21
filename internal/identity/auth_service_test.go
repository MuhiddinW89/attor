package identity

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

type authRepositoryStub struct {
	user           *User
	passwordHash   string
	tokenHash      []byte
	csrfHash       []byte
	sessionExpires time.Time
	session        *Session
}

func (r *authRepositoryStub) FindUserByEmail(_ context.Context, email string) (*User, string, error) {
	if r.user == nil || email != r.user.Email {
		return nil, "", ErrInvalidCredentials
	}
	return r.user, r.passwordHash, nil
}

func (r *authRepositoryStub) CreateSession(_ context.Context, tokenHash, csrfHash []byte, userID uuid.UUID, expiresAt time.Time) error {
	r.tokenHash = append([]byte(nil), tokenHash...)
	r.csrfHash = append([]byte(nil), csrfHash...)
	r.sessionExpires = expiresAt
	r.session = &Session{UserID: userID, StoreID: r.user.StoreID, Role: r.user.Role, CSRFHash: r.csrfHash, ExpiresAt: expiresAt}
	return nil
}

func (r *authRepositoryStub) GetSession(_ context.Context, tokenHash []byte) (*Session, error) {
	if r.session == nil || !bytes.Equal(r.tokenHash, tokenHash) || !r.session.ExpiresAt.After(time.Now()) {
		return nil, ErrInvalidSession
	}
	return r.session, nil
}

func (r *authRepositoryStub) DeleteSession(_ context.Context, tokenHash []byte) error {
	if bytes.Equal(r.tokenHash, tokenHash) {
		r.session = nil
	}
	return nil
}

func TestAuthServiceLoginStoresOnlyDigestsAndRevokesSession(t *testing.T) {
	passwordHash, err := HashPassword("valid-password-123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	repository := &authRepositoryStub{
		user:         &User{ID: uuid.New(), StoreID: uuid.New(), Email: "owner@example.com", Role: RoleOwner},
		passwordHash: passwordHash,
	}
	service := NewAuthService(repository)
	result, err := service.Login(context.Background(), " OWNER@example.com ", "valid-password-123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if bytes.Equal(repository.tokenHash, []byte(result.Token)) || bytes.Equal(repository.csrfHash, []byte(result.CSRFToken)) {
		t.Fatal("raw session or CSRF token was stored")
	}
	session, err := service.Authenticate(context.Background(), result.Token)
	if err != nil || session.UserID != result.User.ID || session.StoreID != result.User.StoreID || session.Role != RoleOwner {
		t.Fatalf("Authenticate() session = %#v, error = %v", session, err)
	}
	if err := service.Logout(context.Background(), result.Token); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if _, err := service.Authenticate(context.Background(), result.Token); err == nil {
		t.Fatal("revoked session authenticated")
	}
}

func TestAuthServiceRejectsInvalidCredentialsAndExpiredSession(t *testing.T) {
	passwordHash, err := HashPassword("valid-password-123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	repository := &authRepositoryStub{
		user:         &User{ID: uuid.New(), StoreID: uuid.New(), Email: "owner@example.com", Role: RoleOwner},
		passwordHash: passwordHash,
	}
	service := NewAuthService(repository)
	if _, err := service.Login(context.Background(), "owner@example.com", "wrong-password"); err != ErrInvalidCredentials {
		t.Fatalf("wrong-password error = %v, want %v", err, ErrInvalidCredentials)
	}
	result, err := service.Login(context.Background(), "owner@example.com", "valid-password-123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	repository.session.ExpiresAt = time.Now().Add(-time.Minute)
	if _, err := service.Authenticate(context.Background(), result.Token); err != ErrInvalidSession {
		t.Fatalf("expired-session error = %v, want %v", err, ErrInvalidSession)
	}
}
