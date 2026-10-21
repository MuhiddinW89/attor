package identity

import (
	"context"
	"testing"
)

type initialOwnerRepositoryStub struct {
	called       bool
	passwordHash string
	user         *User
	err          error
}

func (s *initialOwnerRepositoryStub) CreateInitialOwner(_ context.Context, email, passwordHash, _ string) (*User, error) {
	s.called = true
	s.passwordHash = passwordHash
	if s.user == nil {
		s.user = &User{Email: email, Role: RoleOwner}
	}
	return s.user, s.err
}

func TestCreateInitialOwnerHashesPassword(t *testing.T) {
	repository := &initialOwnerRepositoryStub{}
	service := NewService(repository)
	user, err := service.CreateInitialOwner(context.Background(), "OWNER@example.com", "a-long-test-password", "Attor Store")
	if err != nil {
		t.Fatalf("CreateInitialOwner() error = %v", err)
	}
	if !repository.called || repository.passwordHash == "a-long-test-password" || !VerifyPassword(repository.passwordHash, "a-long-test-password") {
		t.Fatal("owner password was not stored as a verifiable hash")
	}
	if user.Role != RoleOwner || user.Email != "owner@example.com" {
		t.Fatalf("unexpected owner: %#v", user)
	}
}

func TestCreateInitialOwnerRejectsInvalidInput(t *testing.T) {
	repository := &initialOwnerRepositoryStub{}
	service := NewService(repository)
	if _, err := service.CreateInitialOwner(context.Background(), "not-an-email", "short", ""); err == nil {
		t.Fatal("expected invalid input error")
	}
	if repository.called {
		t.Fatal("repository called with invalid owner input")
	}
}

func TestOpaqueTokenIsStoredAsDigest(t *testing.T) {
	token, digest, err := GenerateOpaqueToken()
	if err != nil {
		t.Fatalf("GenerateOpaqueToken() error = %v", err)
	}
	if token == "" || len(digest) != 32 {
		t.Fatalf("unexpected token/digest lengths: %d/%d", len(token), len(digest))
	}
	if token == string(digest) {
		t.Fatal("raw token must not be stored as its digest")
	}
}
