package sales

import (
	"context"
	"strings"
	"testing"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type saleInsertSpy struct {
	query string
	args  []any
}

func (s *saleInsertSpy) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	s.query = query
	s.args = args
	return pgconn.NewCommandTag("INSERT 0 1"), nil
}
func (*saleInsertSpy) Query(context.Context, string, ...any) (pgx.Rows, error) { return nil, nil }
func (*saleInsertSpy) QueryRow(context.Context, string, ...any) pgx.Row        { return nil }

func TestCreateSaleUsesStoreAndCreatorFromSession(t *testing.T) {
	spy := &saleInsertSpy{}
	repository := NewPostgresRepository(spy)
	storeID, userID, clientID := uuid.New(), uuid.New(), uuid.New()
	ctx := identity.WithSession(context.Background(), &identity.Session{
		StoreID: storeID, UserID: userID, Role: identity.RoleEmployee,
	})
	sale := &Sale{ID: uuid.New(), StoreID: uuid.New(), CreatedByUserID: uuid.New(), ClientID: clientID}
	if err := repository.Create(ctx, sale); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !strings.Contains(spy.query, "store_id") || !strings.Contains(spy.query, "created_by_user_id") {
		t.Fatalf("insert does not persist ownership: %s", spy.query)
	}
	if spy.args[1] != storeID || spy.args[2] != userID || spy.args[3] != clientID {
		t.Fatalf("insert ownership values = %#v", spy.args)
	}
	if sale.StoreID != storeID || sale.CreatedByUserID != userID {
		t.Fatal("client-supplied ownership fields were not overwritten")
	}
}
