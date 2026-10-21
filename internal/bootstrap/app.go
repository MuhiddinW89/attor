package bootstrap

import (
	"os"
	"strings"
	"time"

	"github.com/MuhiddinW89/attor/internal/analytics"
	analyticsapp "github.com/MuhiddinW89/attor/internal/application/analytics"
	clientapp "github.com/MuhiddinW89/attor/internal/application/clients"
	reminderapp "github.com/MuhiddinW89/attor/internal/application/reminders"
	saleapp "github.com/MuhiddinW89/attor/internal/application/sales"
	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/MuhiddinW89/attor/internal/identity"
	"github.com/MuhiddinW89/attor/internal/reminders"
	"github.com/MuhiddinW89/attor/internal/sales"
	transport "github.com/MuhiddinW89/attor/internal/transport/http"
	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/MuhiddinW89/attor/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	analyticsRepo := analytics.NewPostgresRepository(db)
	identityRepo := identity.NewPostgresRepository(db)

	// =========================
	// Services
	// =========================

	clientService := clients.NewService(clientRepo)
	saleService := sales.NewService(saleRepo)
	reminderService := reminders.NewService(reminderRepo)
	analyticsService := analytics.NewService(analyticsRepo)
	authService := identity.NewAuthService(identityRepo)

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

	getAnalyticsUC := analyticsapp.NewGetAnalyticsUseCase(
		analyticsService,
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

	analyticsHandler := transport.NewAnalyticsHandler(
		getAnalyticsUC,
	)
	secureCookies := !strings.EqualFold(os.Getenv("APP_ENV"), "development")
	allowedOrigins := transport.AllowedOriginsForEnvironment(secureCookies)
	authHandler := transport.NewAuthHandler(
		authService,
		secureCookies,
		allowedOrigins,
	)
	// =========================
	// Fiber
	// =========================

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			if fiberErr, ok := err.(*fiber.Error); ok && fiberErr.Code < 500 {
				status = fiberErr.Code
			}
			message := "internal server error"
			switch status {
			case fiber.StatusBadRequest:
				message = fiber.ErrBadRequest.Message
			case fiber.StatusNotFound:
				message = "not found"
			case fiber.StatusMethodNotAllowed:
				message = "method not allowed"
			case fiber.StatusRequestEntityTooLarge:
				message = "request too large"
			}
			return c.Status(status).JSON(fiber.Map{"error": message})
		},
	})
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CSRF-Token",
		AllowCredentials: true,
	}))
	app.Use("/api", limiter.New(limiter.Config{
		Max:        60,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "too many requests",
			})
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Attor API is running")
	})
	transport.RegisterAuthRoutes(app, authHandler)

	transport.RegisterRoutes(
		app,
		authHandler,
		clientHandler,
		saleHandler,
		reminderHandler,
		analyticsHandler,
	)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "attor-backend",
		})
	})

	return app, nil
}
