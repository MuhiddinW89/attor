package transport

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/gofiber/fiber/v2"
)

const (
	sessionCookieName = "attor_session"
	csrfCookieName    = "attor_csrf"
	csrfHeaderName    = "X-CSRF-Token"
	sessionLocalKey   = "attor_session_context"
)

var AllowedOrigins = []string{
	"http://localhost:5173",
	"https://attor-frontend.pages.dev",
}

var productionOrigins = []string{"https://attor-frontend.pages.dev"}

type Authenticator interface {
	Login(ctx context.Context, email, password string) (*identity.LoginResult, error)
	Authenticate(ctx context.Context, token string) (*identity.Session, error)
	Logout(ctx context.Context, token string) error
}

type AuthHandler struct {
	service        Authenticator
	secure         bool
	allowedOrigins []string
}

func NewAuthHandler(service Authenticator, secureCookie bool, allowedOrigins ...[]string) *AuthHandler {
	origins := AllowedOrigins
	if secureCookie {
		origins = productionOrigins
	}
	if len(allowedOrigins) > 0 {
		origins = allowedOrigins[0]
	}
	return &AuthHandler{service: service, secure: secureCookie, allowedOrigins: origins}
}

func AllowedOriginsForEnvironment(production bool) []string {
	if production {
		return append([]string(nil), productionOrigins...)
	}
	return append([]string(nil), AllowedOrigins...)
}

func RegisterAuthRoutes(app *fiber.App, handler *AuthHandler) {
	auth := app.Group("/api/v1/auth")
	auth.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderCacheControl, "no-store")
		return c.Next()
	})
	auth.Post("/login", handler.RequireAllowedOrigin, handler.Login)
	auth.Post("/logout", handler.RequireAllowedOrigin, handler.RequireSession, handler.RequireCSRF, handler.Logout)
	auth.Get("/me", handler.RequireSession, handler.Me)
}

func (h *AuthHandler) RequireAllowedOrigin(c *fiber.Ctx) error {
	origin := c.Get(fiber.HeaderOrigin)
	for _, allowed := range h.allowedOrigins {
		if origin == allowed {
			return c.Next()
		}
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "origin not allowed"})
}

func (h *AuthHandler) RequireSession(c *fiber.Ctx) error {
	token := c.Cookies(sessionCookieName)
	session, err := h.service.Authenticate(c.Context(), token)
	if err != nil {
		if errors.Is(err, identity.ErrInvalidSession) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	c.Locals(sessionLocalKey, session)
	c.SetUserContext(identity.WithSession(c.UserContext(), session))
	return c.Next()
}

func (h *AuthHandler) RequireCSRF(c *fiber.Ctx) error {
	session, ok := SessionFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	cookieToken := c.Cookies(csrfCookieName)
	headerToken := c.Get(csrfHeaderName)
	if cookieToken == "" || headerToken == "" ||
		subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 ||
		subtle.ConstantTimeCompare(identity.HashToken(headerToken), session.CSRFHash) != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "csrf validation failed"})
	}
	return c.Next()
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&request); err != nil || len(request.Email) > 320 || request.Password == "" || len(request.Password) > 72 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	result, err := h.service.Login(c.Context(), request.Email, request.Password)
	if errors.Is(err, identity.ErrInvalidCredentials) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	h.setCookies(c, result.Token, result.CSRFToken, result.ExpiresAt)
	return c.JSON(fiber.Map{
		"user_id":    result.User.ID,
		"store_id":   result.User.StoreID,
		"email":      result.User.Email,
		"role":       result.User.Role,
		"csrf_token": result.CSRFToken,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Cookies(sessionCookieName)
	err := h.service.Logout(c.Context(), token)
	h.clearCookies(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	session, ok := SessionFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{
		"user_id":    session.UserID,
		"store_id":   session.StoreID,
		"role":       session.Role,
		"expires_at": session.ExpiresAt,
	})
}

func (h *AuthHandler) setCookies(c *fiber.Ctx, sessionToken, csrfToken string, expiresAt time.Time) {
	sameSite := fiber.CookieSameSiteLaxMode
	if h.secure {
		sameSite = fiber.CookieSameSiteNoneMode
	}
	maxAge := int(time.Until(expiresAt).Seconds())
	c.Cookie(&fiber.Cookie{
		Name: sessionCookieName, Value: sessionToken, Path: "/api/v1", Expires: expiresAt,
		MaxAge: maxAge, HTTPOnly: true, Secure: h.secure, SameSite: sameSite,
	})
	c.Cookie(&fiber.Cookie{
		Name: csrfCookieName, Value: csrfToken, Path: "/api/v1", Expires: expiresAt,
		MaxAge: maxAge, HTTPOnly: false, Secure: h.secure, SameSite: sameSite,
	})
}

func (h *AuthHandler) clearCookies(c *fiber.Ctx) {
	for _, name := range []string{sessionCookieName, csrfCookieName} {
		c.Cookie(&fiber.Cookie{
			Name: name, Value: "", Path: "/api/v1", Expires: time.Unix(1, 0), MaxAge: -1,
			HTTPOnly: name == sessionCookieName, Secure: h.secure,
			SameSite: fiber.CookieSameSiteLaxMode,
		})
	}
}

func SessionFromContext(c *fiber.Ctx) (*identity.Session, bool) {
	session, ok := c.Locals(sessionLocalKey).(*identity.Session)
	return session, ok
}

func RequireRoles(roles ...identity.Role) fiber.Handler {
	allowed := make(map[identity.Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		session, ok := SessionFromContext(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		if _, ok := allowed[session.Role]; !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}
