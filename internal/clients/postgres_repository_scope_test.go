package clients

import (
	"context"
	"strings"
	"testing"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type scopedQuerySpy struct {
	query string
	args  []any
}

func (s *scopedQuerySpy) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (s *scopedQuerySpy) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, nil
}

func (s *scopedQuerySpy) QueryRow(_ context.Context, query string, args ...any) pgx.Row {
	s.query = query
	s.args = args
	return noClientRow{}
}

type noClientRow struct{}

func (noClientRow) Scan(...any) error { return pgx.ErrNoRows }

func TestGetByIDScopesForgedUUIDToSessionStore(t *testing.T) {
	spy := &scopedQuerySpy{}
	repository := NewPostgresRepository(spy)
	storeID := uuid.New()
	userID := uuid.New()
	foreignRecordID := uuid.New()
	ctx := identity.WithSession(context.Background(), &identity.Session{
		StoreID: storeID, UserID: userID, Role: identity.RoleEmployee,
	})

	client, err := repository.GetByID(ctx, foreignRecordID)
	if err != ErrClientNotFound || client != nil {
		t.Fatalf("GetByID() = (%v, %v), want not found for foreign record", client, err)
	}
	if !strings.Contains(spy.query, "id = $1 AND store_id = $2") || !strings.Contains(spy.query, "created_by_user_id = $4") {
		t.Fatalf("query lacks tenant/employee constraints: %s", spy.query)
	}
	if spy.args[0] != foreignRecordID || spy.args[1] != storeID || spy.args[2] != identity.RoleEmployee || spy.args[3] != userID {
		t.Fatalf("query args do not come from requested ID + verified session: %#v", spy.args)
	}
}

func TestGetByIDRejectsMissingSessionBeforeQuery(t *testing.T) {
	repository := NewPostgresRepository(nil)
	if _, err := repository.GetByID(context.Background(), uuid.New()); err != identity.ErrInvalidSession {
		t.Fatalf("GetByID() error = %v, want invalid session", err)
	}
}
