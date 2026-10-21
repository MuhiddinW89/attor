package reminders

import (
	"context"
	"strings"
	"testing"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type reminderScopeSpy struct {
	query string
	args  []any
	rows  int64
}

func (s *reminderScopeSpy) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	s.query = query
	s.args = args
	return pgconn.NewCommandTag("UPDATE 1"), nil
}

func (s *reminderScopeSpy) Query(_ context.Context, query string, args ...any) (pgx.Rows, error) {
	s.query = query
	s.args = args
	return &emptyReminderRows{}, nil
}

func (*reminderScopeSpy) QueryRow(context.Context, string, ...any) pgx.Row {
	return reminderNoRows{}
}

type reminderNoRows struct{}

func (reminderNoRows) Scan(...any) error {
	return pgx.ErrNoRows
}

type emptyReminderRows struct{}

func (*emptyReminderRows) Close()                                       {}
func (*emptyReminderRows) Err() error                                   { return nil }
func (*emptyReminderRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (*emptyReminderRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (*emptyReminderRows) Next() bool                                   { return false }
func (*emptyReminderRows) Scan(...any) error                            { return nil }
func (*emptyReminderRows) Values() ([]any, error)                       { return nil, nil }
func (*emptyReminderRows) RawValues() [][]byte                          { return nil }
func (*emptyReminderRows) Conn() *pgx.Conn                              { return nil }

func reminderTestContext(storeID, userID uuid.UUID) context.Context {
	return identity.WithSession(context.Background(), &identity.Session{
		StoreID: storeID,
		UserID:  userID,
		Role:    identity.RoleEmployee,
	})
}

func TestCreateUsesStoreAndCreatorFromSession(t *testing.T) {
	spy := &reminderScopeSpy{}
	repository := NewPostgresRepository(spy)

	storeID, userID := uuid.New(), uuid.New()
	reminder := &Reminder{
		ID:              uuid.New(),
		StoreID:         uuid.New(),
		CreatedByUserID: uuid.New(),
		ClientID:        uuid.New(),
		SaleID:          uuid.New(),
	}

	if err := repository.Create(reminderTestContext(storeID, userID), reminder); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if !strings.Contains(spy.query, "store_id") ||
		!strings.Contains(spy.query, "created_by_user_id") {
		t.Fatalf("INSERT lacks ownership columns: %s", spy.query)
	}
	if spy.args[1] != storeID || spy.args[2] != userID {
		t.Fatalf("INSERT ownership arguments = %#v", spy.args)
	}
	if reminder.StoreID != storeID || reminder.CreatedByUserID != userID {
		t.Fatal("Create did not overwrite ownership values from the reminder")
	}
}

func TestListPendingScopesToSession(t *testing.T) {
	spy := &reminderScopeSpy{}
	repository := NewPostgresRepository(spy)

	storeID, userID := uuid.New(), uuid.New()
	_, err := repository.ListPending(reminderTestContext(storeID, userID))
	if err != nil {
		t.Fatalf("ListPending() error = %v", err)
	}

	if !strings.Contains(spy.query, "store_id = $1") ||
		!strings.Contains(spy.query, "created_by_user_id = $3") {
		t.Fatalf("query lacks store/employee scope: %s", spy.query)
	}
	if len(spy.args) != 3 ||
		spy.args[0] != storeID ||
		spy.args[1] != identity.RoleEmployee ||
		spy.args[2] != userID {
		t.Fatalf("query arguments do not match session: %#v", spy.args)
	}
}

func TestMarkSentScopesToSession(t *testing.T) {
	spy := &reminderScopeSpy{}
	repository := NewPostgresRepository(spy)

	storeID, userID, reminderID := uuid.New(), uuid.New(), uuid.New()
	err := repository.MarkSent(reminderTestContext(storeID, userID), reminderID)
	if err != nil {
		t.Fatalf("MarkSent() error = %v", err)
	}

	if !strings.Contains(spy.query, "WHERE id = $1 AND store_id = $2") ||
		!strings.Contains(spy.query, "created_by_user_id = $4") {
		t.Fatalf("UPDATE lacks store/employee scope: %s", spy.query)
	}
	if len(spy.args) != 4 ||
		spy.args[0] != reminderID ||
		spy.args[1] != storeID ||
		spy.args[2] != identity.RoleEmployee ||
		spy.args[3] != userID {
		t.Fatalf("UPDATE arguments do not match request and session: %#v", spy.args)
	}
}

func TestListPendingRejectsMissingSession(t *testing.T) {
	repository := NewPostgresRepository(nil)
	_, err := repository.ListPending(context.Background())
	if err != identity.ErrInvalidSession {
		t.Fatalf("ListPending() error = %v, want invalid session", err)
	}
}

func TestMarkSentRejectsMissingSession(t *testing.T) {
	repository := NewPostgresRepository(nil)
	err := repository.MarkSent(context.Background(), uuid.New())
	if err != identity.ErrInvalidSession {
		t.Fatalf("MarkSent() error = %v, want invalid session", err)
	}
}
