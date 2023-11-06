package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"sambrlon/internal/tokens"
	"sambrlon/internal/tokens/usecase"
)

type Handler struct {
	tokenUC tokens.UseCase
}

func NewHandler(uc tokens.UseCase) *Handler {
	return &Handler{
		tokenUC: uc,
	}
}

func (h *Handler) AddTicker(c *fiber.Ctx) error {
	var params tokens.AddTickerParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	err := h.tokenUC.AddTicker(params.Ticker)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrFetchTicker):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, usecase.ErrSaveTicker):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unknown"})
		}
	}

	h.tokenUC.RunFetch(params.Ticker)
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) FetchTicker(c *fiber.Ctx) error {
	var params tokens.FetchTickerParams
	if err := c.QueryParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	tickerData, err := h.tokenUC.FetchTickerData(params.Ticker, params.DateFrom, params.DateTo)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrFetchTicker):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unknown"})
		}
	}

	if len(tickerData) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	return c.Status(fiber.StatusOK).JSON(tickerData)
}
