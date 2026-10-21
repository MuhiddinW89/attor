package analytics

import (
	"strings"
	"testing"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/google/uuid"
)

func TestBuildSalesFilterScopesStoreAndEmployee(t *testing.T) {
	storeID := uuid.New()
	userID := uuid.New()
	where, args := buildSalesFilter(Filters{}, &identity.Session{
		StoreID: storeID, UserID: userID, Role: identity.RoleEmployee,
	})
	if !strings.Contains(where, "sales.store_id = $1") || !strings.Contains(where, "sales.created_by_user_id = $3") {
		t.Fatalf("analytics filter lacks store or creator scope: %s", where)
	}
	if len(args) != 3 || args[0] != storeID || args[1] != identity.RoleEmployee || args[2] != userID {
		t.Fatalf("analytics scope args = %#v", args)
	}
}
