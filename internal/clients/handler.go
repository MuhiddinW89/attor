package clients

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(
	service Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(
	c *fiber.Ctx,
) error {

	var req CreateClientRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid request body",
			},
		)
	}

	client, err := h.service.Create(
		c.Context(),
		CreateClientInput{
			FullName:  req.FullName,
			Phone:     req.Phone,
			Instagram: req.Instagram,
		},
	)

	if err != nil {

		if errors.Is(err, ErrClientAlreadyExists) {
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
		ClientResponse{
			ID:        client.ID.String(),
			FullName:  client.FullName,
			Phone:     client.Phone,
			Instagram: client.Instagram,
		},
	)
}

func (h *Handler) List(
	c *fiber.Ctx,
) error {

	clients, err := h.service.List(
		c.Context(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	response := make(
		[]ClientListItem,
		0,
		len(clients),
	)

	for _, client := range clients {

		response = append(
			response,
			ClientListItem{
				ID:       client.ID.String(),
				FullName: client.FullName,
				Phone:    client.Phone,
			},
		)
	}

	return c.JSON(response)
}

func (h *Handler) GetByID(
	c *fiber.Ctx,
) error {

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

		if errors.Is(err, ErrClientNotFound) {
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
		ClientDetailsResponse{
			ID:        client.ID.String(),
			FullName:  client.FullName,
			Phone:     client.Phone,
			Instagram: client.Instagram,
			BirthDate: birthDate,
		},
	)
}
