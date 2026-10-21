package transport

import (
	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	authHandler *AuthHandler,
	clientHandler *ClientHandler,
	saleHandler *SaleHandler,
	reminderHandler *ReminderHandler,
	analyticsHandler *AnalyticsHandler,
) {

	v1 := app.Group(
		"/api/v1",
		authHandler.RequireSession,
		RequireRoles(identity.RoleOwner, identity.RoleEmployee),
	)

	// Clients

	v1.Post("/clients", authHandler.RequireAllowedOrigin, authHandler.RequireCSRF, clientHandler.Create)
	v1.Get("/clients", clientHandler.List)
	v1.Get("/clients/summary", clientHandler.Summary)
	v1.Get("/clients/:id", clientHandler.GetByID)

	// Sales

	v1.Post("/sales", authHandler.RequireAllowedOrigin, authHandler.RequireCSRF, saleHandler.Create)

	v1.Get(
		"/clients/:id/history",
		clientHandler.History,
	)

	v1.Patch(
		"/reminders/:id/sent",
		authHandler.RequireAllowedOrigin,
		authHandler.RequireCSRF,
		reminderHandler.MarkSent,
	)

	v1.Get(
		"/reminders",
		reminderHandler.List,
	)

	// Analytics

	v1.Get(
		"/analytics",
		analyticsHandler.Get,
	)
}
