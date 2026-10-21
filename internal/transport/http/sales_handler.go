package transport

import (
	"strings"
	"time"

	saleapp "github.com/MuhiddinW89/attor/internal/application/sales"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SaleHandler struct {
	useCase *saleapp.CreateSaleUseCase
}

func NewSaleHandler(
	useCase *saleapp.CreateSaleUseCase,
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
	if (req.ClientID == "" && (!validName(req.FullName, 255) || !validPhone(req.Phone))) ||
		!validName(req.PerfumeName, 255) || req.VolumeML <= 0 || req.VolumeML > maxVolumeML ||
		!validSalePrice(req.Price) || (req.Comment != nil && runeLength(*req.Comment) > 5000) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid sale fields",
		})
	}
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.PerfumeName = strings.TrimSpace(req.PerfumeName)
	if req.Comment != nil {
		comment := strings.TrimSpace(*req.Comment)
		req.Comment = &comment
	}
	var clientID *uuid.UUID
	if req.ClientID != "" {
		parsed, err := uuid.Parse(req.ClientID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid client id"})
		}
		clientID = &parsed
	}

	var birthDate *time.Time

	if req.Birthday != "" {
		parsed, err := time.Parse("02.01.2006", req.Birthday)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"error": "invalid birthday format, expected DD.MM.YYYY",
				},
			)
		}

		birthDate = &parsed
		if parsed.After(time.Now()) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "birthday cannot be in the future",
			})
		}
	}

	sale, err := h.useCase.Execute(
		c.UserContext(),
		saleapp.CreateSaleRequest{
			ClientID:    clientID,
			Phone:       req.Phone,
			FullName:    req.FullName,
			BirthDate:   birthDate,
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
				"error": "internal server error",
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
