package http

import (
	"github.com/gofiber/fiber/v2"
	"sambrlon/internal/middleware"
	"sambrlon/internal/tokens"
	"sambrlon/internal/tokens/repository"
)

func RegisterRoutes(app *fiber.App, handler tokens.Handler, db *repository.PostgresDB) {
	api := app.Group("/api")

	api.Use(middleware.AuthMiddleware(db))

	api.Post("/add_ticker", handler.AddTicker)
	api.Get("/fetch", handler.FetchTicker)
}
