package transport

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterRoutesRequiresSession(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app, NewAuthHandler(&authServiceStub{}, false), nil, nil, nil, nil)

	request := httptest.NewRequest("GET", "/api/v1/clients", nil)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusUnauthorized)
	}
}
