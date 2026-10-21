package transport

import (
	"errors"
	"strings"

	clientapp "github.com/MuhiddinW89/attor/internal/application/clients"
	"github.com/MuhiddinW89/attor/internal/clients"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClientHandler struct {
	service        clients.Service
	historyUseCase *clientapp.ClientHistoryUseCase
}

func NewClientHandler(
	service clients.Service,
	historyUseCase *clientapp.ClientHistoryUseCase,
) *ClientHandler {

	return &ClientHandler{
		service:        service,
		historyUseCase: historyUseCase,
	}
}

func (h *ClientHandler) Create(c *fiber.Ctx) error {
	var req clients.CreateClientRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid request body",
			},
		)
	}
	if !validName(req.FullName, 255) || !validPhone(req.Phone) ||
		(req.Instagram != nil && runeLength(*req.Instagram) > 255) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid client fields",
		})
	}
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Instagram != nil {
		instagram := strings.TrimSpace(*req.Instagram)
		req.Instagram = &instagram
	}

	client, err := h.service.Create(
		c.UserContext(),
		clients.CreateClientInput{
			FullName:  req.FullName,
			Phone:     req.Phone,
			Instagram: req.Instagram,
		},
	)

	if err != nil {

		if errors.Is(err, clients.ErrClientAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(
				fiber.Map{
					"error": "client already exists",
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		clients.ClientResponse{
			ID:        client.ID.String(),
			FullName:  client.FullName,
			Phone:     client.Phone,
			Instagram: client.Instagram,
		},
	)
}

func (h *ClientHandler) List(c *fiber.Ctx) error {

	search := c.Query("search")
	if runeLength(search) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "search is too long"})
	}

	clientList, err := h.service.List(
		c.UserContext(),
		search,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	response := make(
		[]clients.ClientListItem,
		0,
		len(clientList),
	)

	for _, client := range clientList {
		response = append(
			response,
			clients.ClientListItem{
				ID:       client.ID.String(),
				FullName: client.FullName,
				Phone:    client.Phone,
			},
		)
	}

	return c.JSON(response)
}

func (h *ClientHandler) GetByID(c *fiber.Ctx) error {

	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid client id",
			},
		)
	}

	client, err := h.service.GetByID(
		c.UserContext(),
		id,
	)

	if err != nil {

		if errors.Is(err, clients.ErrClientNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": "client not found",
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	var birthDate *string

	if client.BirthDate != nil {
		formatted := client.BirthDate.Format("2006-01-02")
		birthDate = &formatted
	}

	return c.JSON(
		clients.ClientDetailsResponse{
			ID:        client.ID.String(),
			FullName:  client.FullName,
			Phone:     client.Phone,
			Instagram: client.Instagram,
			BirthDate: birthDate,
		},
	)
}

func (h *ClientHandler) History(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid client id",
			},
		)
	}

	history, err := h.historyUseCase.Execute(
		c.UserContext(),
		id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	return c.JSON(history)
}

func (h *ClientHandler) Summary(c *fiber.Ctx) error {

	search := c.Query("search")
	if runeLength(search) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "search is too long"})
	}

	clientList, err := h.service.ListWithStats(
		c.UserContext(),
		search,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	response := make(
		[]clients.ClientListStatsItem,
		0,
		len(clientList),
	)

	for _, client := range clientList {
		response = append(
			response,
			clients.ClientListStatsItem{
				ID:          client.ID,
				FullName:    client.FullName,
				Phone:       client.Phone,
				Instagram:   client.Instagram,
				TotalSales:  client.TotalSales,
				TotalAmount: client.TotalAmount,
			},
		)
	}

	return c.JSON(response)
}
