package bootstrap

import (
	clientapp "github.com/MuhiddinW89/attor/internal/application/clients"
	reminderapp "github.com/MuhiddinW89/attor/internal/application/reminders"
	saleapp "github.com/MuhiddinW89/attor/internal/application/sales"
	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/MuhiddinW89/attor/internal/sales"
	transport "github.com/MuhiddinW89/attor/internal/transport/http"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func New(
	cfg *config.Config,
) (*fiber.App, error) {

	// =========================
	// Config
	// =========================

	// =========================
	// Database
	// =========================

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

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

	createSaleUC := saleapp.NewCreateSaleUseCase(
		clientService,
		saleService,
		reminderService,
	)

	clientHistoryUC := clientapp.NewClientHistoryUseCase(
		clientService,
		saleService,
		reminderService,
	)

	markReminderSentUC := reminderapp.NewMarkReminderSentUseCase(
		reminderService,
	)

	listRemindersUC := reminderapp.NewListUseCase(
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
		listRemindersUC,
	)
	// =========================
	// Fiber
	// =========================

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Attor API is running")
	})

	transport.RegisterRoutes(
		app,
		clientHandler,
		saleHandler,
		reminderHandler,
	)

	return app, nil
}
