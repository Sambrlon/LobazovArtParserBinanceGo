package main

import (
	"log"
	"sambrlon/config"
	"sambrlon/internal/httpServer"
	"sambrlon/internal/tokens/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	server := httpServer.NewServer(cfg, db)
	server.Run()
}
