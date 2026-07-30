package main

import (
	"log"

	"github.com/MuhiddinW89/attor/internal/clients"
	transport "github.com/MuhiddinW89/attor/internal/transport/http"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// =========================
	// Clients
	// =========================

	clientRepo := clients.NewPostgresRepository(db)
	clientService := clients.NewService(clientRepo)
	clientHandler := transport.NewClientHandler(clientService)

	// =========================
	// HTTP Server
	// =========================

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Attor API is running")
	})

	transport.RegisterRoutes(app, clientHandler)

	log.Println("Server started on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
