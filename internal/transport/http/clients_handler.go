package transport

import (
	"errors"

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

	client, err := h.service.Create(
		c.Context(),
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
					"error": err.Error(),
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
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

	clientList, err := h.service.List(
		c.Context(),
		search,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
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
		c.Context(),
		id,
	)

	if err != nil {

		if errors.Is(err, clients.ErrClientNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
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
		c.Context(),
		id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(history)
}
