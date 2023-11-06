package httpServer

import (
	"github.com/gofiber/fiber/v2"
	"sambrlon/config"
	"sambrlon/internal/tokens/repository"
)

type Server struct {
	app    *fiber.App
	config config.Config
}

func NewServer(cfg config.Config, db *repository.PostgresDB) *Server {
	app := fiber.New()

	// app.Use(middleware.AuthMiddleware(&cfg, db))

	return &Server{app: app, config: cfg}
}

func (s *Server) Run() {
	MapHandler(s.app, s.config)

	err := s.app.Listen(":" + s.config.Port)
	if err != nil {
		panic(err)
	}
}
