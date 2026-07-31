package transport

import (
	"time"

	"github.com/MuhiddinW89/attor/internal/application"
	"github.com/gofiber/fiber/v2"
)

type SaleHandler struct {
	useCase *application.CreateSaleUseCase
}

func NewSaleHandler(
	useCase *application.CreateSaleUseCase,
) *SaleHandler {
	return &SaleHandler{
		useCase: useCase,
	}
}


func (h *SaleHandler) Create(c *fiber.Ctx) error {
	var req CreateSaleRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid request body",
			},
		)
	}

	sale, err := h.useCase.Execute(
		c.Context(),
		application.CreateSaleRequest{
			Phone:       req.Phone,
			FullName:    req.FullName,
			PerfumeName: req.PerfumeName,
			VolumeML:    req.VolumeML,
			Price:       req.Price,
			Comment:     req.Comment,
			SaleDate:    time.Now(),
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"id":           sale.ID,
			"client_id":    sale.ClientID,
			"perfume_name": sale.PerfumeName,
			"volume_ml":    sale.VolumeML,
			"price":        sale.Price,
			"comment":      sale.Comment,
			"sale_date":    sale.SaleDate,
			"created_at":   sale.CreatedAt,
			"updated_at":   sale.UpdatedAt,
		},
	)
}
