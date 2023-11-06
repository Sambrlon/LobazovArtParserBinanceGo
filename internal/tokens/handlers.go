package tokens

import "github.com/gofiber/fiber/v2"

type Handler interface {
	AddTicker(c *fiber.Ctx) error
	FetchTicker(c *fiber.Ctx) error
}
