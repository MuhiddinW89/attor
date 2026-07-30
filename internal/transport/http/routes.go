package transport

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, clientHandler *ClientHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/clients", clientHandler.Create)
	v1.Get("/clients", clientHandler.List)
	v1.Get("/clients/:id", clientHandler.GetByID)
}
