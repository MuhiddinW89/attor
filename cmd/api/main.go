package main

import (
	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := clients.NewPostgresRepository(db)

	service := clients.NewService(repo)

	handler := clients.NewHandler(service)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Attor API is running")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post(
		"/clients",
		handler.Create,
	)

	v1.Get(
		"/clients",
		handler.List,
	)

	v1.Get(
		"/clients/:id",
		handler.GetByID,
	)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
