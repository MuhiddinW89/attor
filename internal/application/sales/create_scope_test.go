package sales

import (
	"context"
	"testing"

	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/MuhiddinW89/attor/internal/reminders"
	domainSales "github.com/MuhiddinW89/attor/internal/sales"
	"github.com/google/uuid"
)

type clientScopeStub struct {
	requestedID uuid.UUID
	session     *identity.Session
	getErr      error
}

func (s *clientScopeStub) Create(context.Context, clients.CreateClientInput) (*clients.Client, error) {
	return nil, nil
}
func (s *clientScopeStub) List(context.Context, string) ([]*clients.Client, error) {
	return nil, nil
}
func (s *clientScopeStub) ListWithStats(context.Context, string) ([]*clients.ClientListStatsItem, error) {
	return nil, nil
}
func (s *clientScopeStub) GetByID(ctx context.Context, id uuid.UUID) (*clients.Client, error) {
	s.requestedID = id
	s.session, _ = identity.SessionFromContext(ctx)
	return nil, s.getErr
}
func (s *clientScopeStub) GetByPhone(context.Context, string) (*clients.Client, error) {
	return nil, clients.ErrClientNotFound
}

type saleScopeStub struct{ createCalled bool }

func (s *saleScopeStub) Create(context.Context, domainSales.CreateSaleInput) (*domainSales.Sale, error) {
	s.createCalled = true
	return nil, nil
}
func (s *saleScopeStub) GetByID(context.Context, uuid.UUID) (*domainSales.Sale, error) {
	return nil, nil
}
func (s *saleScopeStub) ListByClientID(context.Context, uuid.UUID) ([]*domainSales.Sale, error) {
	return nil, nil
}

type reminderScopeStub struct{}

func (reminderScopeStub) Create(context.Context, reminders.CreateReminderInput) (*reminders.Reminder, error) {
	return nil, nil
}
func (reminderScopeStub) ListPending(context.Context) ([]*reminders.Reminder, error) {
	return nil, nil
}
func (reminderScopeStub) ListDetailed(context.Context, string) ([]*reminders.ReminderListItem, error) {
	return nil, nil
}
func (reminderScopeStub) MarkSent(context.Context, uuid.UUID) error { return nil }
func (reminderScopeStub) GetNearestByClientID(context.Context, uuid.UUID) (*reminders.Reminder, error) {
	return nil, nil
}
func (reminderScopeStub) ListPendingDetailed(context.Context) ([]*reminders.ReminderListItem, error) {
	return nil, nil
}

func TestCreateSaleWithForeignClientIDStopsBeforeSale(t *testing.T) {
	storeID := uuid.New()
	userID := uuid.New()
	forgedClientID := uuid.New()
	clientService := &clientScopeStub{getErr: clients.ErrClientNotFound}
	saleService := &saleScopeStub{}
	useCase := NewCreateSaleUseCase(clientService, saleService, reminderScopeStub{})
	ctx := identity.WithSession(context.Background(), &identity.Session{
		StoreID: storeID, UserID: userID, Role: identity.RoleEmployee,
	})
	_, err := useCase.Execute(ctx, CreateSaleRequest{ClientID: &forgedClientID})
	if err != clients.ErrClientNotFound {
		t.Fatalf("Execute() error = %v, want scoped client not found", err)
	}
	if clientService.requestedID != forgedClientID || clientService.session == nil ||
		clientService.session.StoreID != storeID || clientService.session.UserID != userID {
		t.Fatal("client lookup did not receive the verified session scope")
	}
	if saleService.createCalled {
		t.Fatal("sale was created using an inaccessible client ID")
	}
}
