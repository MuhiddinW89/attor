package transport

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(
	app *fiber.App,
	clientHandler *ClientHandler,
	saleHandler *SaleHandler,
) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Clients

	v1.Post("/clients", clientHandler.Create)
	v1.Get("/clients", clientHandler.List)
	v1.Get("/clients/:id", clientHandler.GetByID)

	// Sales

	v1.Post("/sales", saleHandler.Create)

	v1.Get(
		"/clients/:id/history",
		clientHandler.History,
	)
}
