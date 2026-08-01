package main

import (
	"log"

	"github.com/MuhiddinW89/attor/internal/application"
	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/MuhiddinW89/attor/internal/sales"
	transport "github.com/MuhiddinW89/attor/internal/transport/http"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// =========================
	// Config
	// =========================

	cfg := config.Load()

	// =========================
	// Database
	// =========================

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	// =========================
	// Repositories
	// =========================

	clientRepo := clients.NewPostgresRepository(db)
	saleRepo := sales.NewPostgresRepository(db)
	reminderRepo := reminders.NewPostgresRepository(db)

	// =========================
	// Services
	// =========================

	clientService := clients.NewService(clientRepo)
	saleService := sales.NewService(saleRepo)
	reminderService := reminders.NewService(reminderRepo)

	// =========================
	// Use Cases
	// =========================

	createSaleUC := application.NewCreateSaleUseCase(
		clientService,
		saleService,
		reminderService,
	)

	clientHistoryUC := application.NewClientHistoryUseCase(
		clientService,
		saleService,
		reminderService,
	)

	markReminderSentUC := application.NewMarkReminderSentUseCase(
		reminderService,
	)

	// =========================
	// Handlers
	// =========================

	clientHandler := transport.NewClientHandler(
		clientService,
		clientHistoryUC,
	)

	saleHandler := transport.NewSaleHandler(
		createSaleUC,
	)

	reminderHandler := transport.NewReminderHandler(
		markReminderSentUC,
	)
	// =========================
	// Fiber
	// =========================

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Attor API is running")
	})

	transport.RegisterRoutes(
		app,
		clientHandler,
		saleHandler,
		reminderHandler,
	)

	log.Println("Server started on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
