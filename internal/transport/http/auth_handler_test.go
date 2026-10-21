package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	testSessionToken = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
	testCSRFToken    = "csrf-token-with-enough-length-123456789"
)

type authServiceStub struct {
	loginErr  error
	session   *identity.Session
	loggedOut bool
}

func (s *authServiceStub) Login(_ context.Context, email, password string) (*identity.LoginResult, error) {
	if s.loginErr != nil {
		return nil, s.loginErr
	}
	if email != "owner@example.com" || password != "correct-password" {
		return nil, identity.ErrInvalidCredentials
	}
	user := &identity.User{
		ID:      uuid.MustParse("10000000-0000-0000-0000-000000000001"),
		StoreID: uuid.MustParse("20000000-0000-0000-0000-000000000002"),
		Email:   email, Role: identity.RoleOwner,
	}
	s.session = &identity.Session{
		UserID: user.ID, StoreID: user.StoreID, Role: user.Role,
		CSRFHash: identity.HashToken(testCSRFToken), ExpiresAt: time.Now().Add(time.Hour),
	}
	return &identity.LoginResult{
		User: user, Token: testSessionToken, CSRFToken: testCSRFToken,
		ExpiresAt: s.session.ExpiresAt,
	}, nil
}

func (s *authServiceStub) Authenticate(_ context.Context, token string) (*identity.Session, error) {
	if token != testSessionToken || s.session == nil || !s.session.ExpiresAt.After(time.Now()) {
		return nil, identity.ErrInvalidSession
	}
	return s.session, nil
}

func (s *authServiceStub) Logout(_ context.Context, token string) error {
	if token != testSessionToken {
		return identity.ErrInvalidSession
	}
	s.loggedOut = true
	s.session = nil
	return nil
}

func newAuthTestApp(service Authenticator) *fiber.App {
	app := fiber.New()
	RegisterAuthRoutes(app, NewAuthHandler(service, false))
	return app
}

func TestAuthLoginMeLogoutFlow(t *testing.T) {
	service := &authServiceStub{}
	app := newAuthTestApp(service)

	loginRequest := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"owner@example.com","password":"correct-password"}`))
	loginRequest.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	loginRequest.Header.Set(fiber.HeaderOrigin, AllowedOrigins[0])
	loginResponse, err := app.Test(loginRequest)
	if err != nil {
		t.Fatalf("login request error = %v", err)
	}
	defer loginResponse.Body.Close()
	if loginResponse.StatusCode != fiber.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResponse.StatusCode, fiber.StatusOK)
	}
	var loginBody map[string]interface{}
	if err := json.NewDecoder(loginResponse.Body).Decode(&loginBody); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loginBody["csrf_token"] != testCSRFToken {
		t.Fatal("login response missing CSRF token")
	}
	var sessionCookie string
	for _, cookie := range loginResponse.Cookies() {
		if cookie.Name == sessionCookieName {
			sessionCookie = cookie.Value
			if !cookie.HttpOnly || cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
				t.Fatal("development session cookie has unexpected security attributes")
			}
		}
	}
	if sessionCookie != testSessionToken {
		t.Fatal("session cookie missing")
	}

	meRequest := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	meRequest.AddCookie(&http.Cookie{Name: sessionCookieName, Value: testSessionToken})
	meResponse, err := app.Test(meRequest)
	if err != nil {
		t.Fatalf("me request error = %v", err)
	}
	defer meResponse.Body.Close()
	if meResponse.StatusCode != fiber.StatusOK {
		t.Fatalf("me status = %d, want %d", meResponse.StatusCode, fiber.StatusOK)
	}

	logoutRequest := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	logoutRequest.AddCookie(&http.Cookie{Name: sessionCookieName, Value: testSessionToken})
	logoutRequest.AddCookie(&http.Cookie{Name: csrfCookieName, Value: testCSRFToken})
	logoutRequest.Header.Set(fiber.HeaderOrigin, AllowedOrigins[0])
	logoutRequest.Header.Set(csrfHeaderName, testCSRFToken)
	logoutResponse, err := app.Test(logoutRequest)
	if err != nil {
		t.Fatalf("logout request error = %v", err)
	}
	defer logoutResponse.Body.Close()
	if logoutResponse.StatusCode != fiber.StatusNoContent || !service.loggedOut {
		t.Fatalf("logout status = %d, loggedOut = %t", logoutResponse.StatusCode, service.loggedOut)
	}
}

func TestLoginRemainsPublicWhenBusinessRoutesAreRegistered(t *testing.T) {
	service := &authServiceStub{}
	handler := NewAuthHandler(service, false)
	app := fiber.New()
	RegisterAuthRoutes(app, handler)
	RegisterRoutes(app, handler, nil, nil, nil, nil)

	request := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"owner@example.com","password":"correct-password"}`))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	request.Header.Set(fiber.HeaderOrigin, AllowedOrigins[0])
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("login request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("login status = %d, want %d", response.StatusCode, fiber.StatusOK)
	}
}

func TestAuthRejectsBadPasswordAndMissingSession(t *testing.T) {
	app := newAuthTestApp(&authServiceStub{})
	request := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"owner@example.com","password":"wrong"}`))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	request.Header.Set(fiber.HeaderOrigin, AllowedOrigins[0])
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("login request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("bad-password status = %d, want %d", response.StatusCode, fiber.StatusUnauthorized)
	}

	unauthorizedResponse, err := app.Test(httptest.NewRequest("GET", "/api/v1/auth/me", nil))
	if err != nil {
		t.Fatalf("unauthorized me request error = %v", err)
	}
	defer unauthorizedResponse.Body.Close()
	if unauthorizedResponse.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("unauthorized me status = %d, want %d", unauthorizedResponse.StatusCode, fiber.StatusUnauthorized)
	}
}

func TestAuthRejectsExpiredSession(t *testing.T) {
	service := &authServiceStub{session: &identity.Session{ExpiresAt: time.Now().Add(-time.Minute)}}
	app := newAuthTestApp(service)
	request := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: testSessionToken})
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expired-session request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expired-session status = %d, want %d", response.StatusCode, fiber.StatusUnauthorized)
	}
}

func TestAuthLogoutRequiresCSRFAndAllowedOrigin(t *testing.T) {
	service := &authServiceStub{}
	service.session = &identity.Session{
		UserID: uuid.New(), StoreID: uuid.New(), Role: identity.RoleOwner,
		CSRFHash: identity.HashToken(testCSRFToken), ExpiresAt: time.Now().Add(time.Hour),
	}
	app := newAuthTestApp(service)
	request := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: testSessionToken})
	request.AddCookie(&http.Cookie{Name: csrfCookieName, Value: testCSRFToken})
	request.Header.Set(fiber.HeaderOrigin, "https://evil.example")
	request.Header.Set(csrfHeaderName, testCSRFToken)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("logout request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusForbidden || service.loggedOut {
		t.Fatalf("logout status = %d, loggedOut = %t", response.StatusCode, service.loggedOut)
	}
}

func TestProductionCookiesAndCSRFRejection(t *testing.T) {
	service := &authServiceStub{}
	app := fiber.New()
	RegisterAuthRoutes(app, NewAuthHandler(service, true))
	loginRequest := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"owner@example.com","password":"correct-password"}`))
	loginRequest.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	loginRequest.Header.Set(fiber.HeaderOrigin, productionOrigins[0])
	loginResponse, err := app.Test(loginRequest)
	if err != nil {
		t.Fatalf("login request error = %v", err)
	}
	defer loginResponse.Body.Close()
	if loginResponse.StatusCode != fiber.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResponse.StatusCode, fiber.StatusOK)
	}
	for _, cookie := range loginResponse.Cookies() {
		if !cookie.Secure || cookie.SameSite != http.SameSiteNoneMode {
			t.Fatalf("production cookie has insecure attributes: %#v", cookie)
		}
		if cookie.Name == sessionCookieName && !cookie.HttpOnly {
			t.Fatal("production session cookie is not HttpOnly")
		}
	}

	request := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: testSessionToken})
	request.AddCookie(&http.Cookie{Name: csrfCookieName, Value: testCSRFToken})
	request.Header.Set(fiber.HeaderOrigin, productionOrigins[0])
	request.Header.Set(csrfHeaderName, "incorrect-csrf-token")
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("logout request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusForbidden || service.loggedOut {
		t.Fatalf("csrf-rejected logout status = %d, loggedOut = %t", response.StatusCode, service.loggedOut)
	}
}

func TestProductionAuthRejectsDevelopmentOrigin(t *testing.T) {
	app := fiber.New()
	RegisterAuthRoutes(app, NewAuthHandler(&authServiceStub{}, true))
	request := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"owner@example.com","password":"correct-password"}`))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	request.Header.Set(fiber.HeaderOrigin, "http://localhost:5173")
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("login request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusForbidden {
		t.Fatalf("development-origin status = %d, want %d", response.StatusCode, fiber.StatusForbidden)
	}
}

func TestRequireRolesRejectsEmployee(t *testing.T) {
	app := fiber.New()
	app.Get("/owner", func(c *fiber.Ctx) error {
		c.Locals(sessionLocalKey, &identity.Session{Role: identity.RoleEmployee})
		return c.Next()
	}, RequireRoles(identity.RoleOwner), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})
	response, err := app.Test(httptest.NewRequest("GET", "/owner", nil))
	if err != nil {
		t.Fatalf("owner route request error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusForbidden {
		t.Fatalf("employee owner-route status = %d, want %d", response.StatusCode, fiber.StatusForbidden)
	}
}
