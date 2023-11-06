package httpServer

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"sambrlon/config"
	httpToken "sambrlon/internal/tokens/delivery/http"
	"sambrlon/internal/tokens/repository"
	"sambrlon/internal/tokens/usecase"
	"sambrlon/pkg/binance"
)

func MapHandler(app *fiber.App, cfg config.Config) {
	dbInstance, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Error creating PostgreSQL connection: %v", err)
	}

	binanceClient := binance.NewBinanceClient(cfg.Binance.APIKey, cfg.Binance.SecretKey)
	tickerUsecase := usecase.NewUsecase(dbInstance, binanceClient)

	tokenHandler := httpToken.NewHandler(tickerUsecase)

	httpToken.RegisterRoutes(app, tokenHandler, dbInstance)

	log.Println("Успешно создан экземпляр tickerUsecase.")
}
