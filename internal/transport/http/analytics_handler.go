package transport

import (
	"strconv"

	"github.com/MuhiddinW89/attor/internal/analytics"
	analyticsapp "github.com/MuhiddinW89/attor/internal/application/analytics"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	useCase *analyticsapp.GetAnalyticsUseCase
}

func NewAnalyticsHandler(
	useCase *analyticsapp.GetAnalyticsUseCase,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		useCase: useCase,
	}
}

func (h *AnalyticsHandler) Get(c *fiber.Ctx) error {
	filters, err := parseAnalyticsFilters(c)
	if err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			fiber.Map{
				"error": "invalid analytics filters",
			},
		)
	}

	result, err := h.useCase.Execute(
		c.UserContext(),
		filters,
	)

	if err != nil {
		return c.Status(
			fiber.StatusInternalServerError,
		).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}

	return c.JSON(result)
}

func parseAnalyticsFilters(
	c *fiber.Ctx,
) (analytics.Filters, error) {

	var filters analytics.Filters

	// =========================================================
	// Year
	// =========================================================

	if value := c.Query("year"); value != "" {
		number, err := strconv.Atoi(value)

		if err != nil {
			return filters, fiber.ErrBadRequest
		}

		if number < 2000 || number > 2100 {
			return filters, fiber.ErrBadRequest
		}

		filters.Year = &number
	}

	// =========================================================
	// Month
	//
	// Requires Year.
	//
	// 1 - 12
	// =========================================================

	if value := c.Query("month"); value != "" {
		if filters.Year == nil {
			return filters, fiber.ErrBadRequest
		}

		number, err := strconv.Atoi(value)

		if err != nil || number < 1 || number > 12 {
			return filters, fiber.ErrBadRequest
		}

		filters.Month = &number
	}

	// =========================================================
	// Week
	//
	// Week is a calendar week INSIDE the selected month.
	//
	// 01 - 06
	//
	// Requires Year + Month.
	// =========================================================

	if value := c.Query("week"); value != "" {
		if filters.Year == nil || filters.Month == nil {
			return filters, fiber.ErrBadRequest
		}

		number, err := strconv.Atoi(value)

		if err != nil || number < 1 || number > 6 {
			return filters, fiber.ErrBadRequest
		}

		filters.Week = &number
	}

	// =========================================================
	// Day
	//
	// Requires Year + Month.
	//
	// 01 - 31
	// =========================================================

	if value := c.Query("day"); value != "" {
		if filters.Year == nil || filters.Month == nil {
			return filters, fiber.ErrBadRequest
		}

		number, err := strconv.Atoi(value)

		if err != nil || number < 1 || number > 31 {
			return filters, fiber.ErrBadRequest
		}

		filters.Day = &number
	}

	return filters, nil
}
