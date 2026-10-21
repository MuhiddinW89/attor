package identity

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner    Role = "owner"
	RoleEmployee Role = "employee"
)

type User struct {
	ID        uuid.UUID
	StoreID   uuid.UUID
	Email     string
	Role      Role
	CreatedAt time.Time
}

type Store struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

type Session struct {
	UserID    uuid.UUID
	StoreID   uuid.UUID
	Role      Role
	CSRFHash  []byte
	ExpiresAt time.Time
}

type sessionContextKey struct{}

func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey{}, session)
}

func SessionFromContext(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(*Session)
	return session, ok && session != nil
}

type EmployeeInvite struct {
	StoreID         uuid.UUID
	Email           string
	InvitedByUserID uuid.UUID
	ExpiresAt       time.Time
}
