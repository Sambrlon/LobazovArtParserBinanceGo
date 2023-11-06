package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB      DBConfig
	Binance BinanceConfig
	Port    string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type BinanceConfig struct {
	APIKey    string
	SecretKey string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Dbname:   os.Getenv("DB_NAME"),
		},
		Binance: BinanceConfig{
			APIKey:    os.Getenv("API_KEY"),
			SecretKey: os.Getenv("SECRET_KEY"),
		},
		Port: os.Getenv("PORT"),
	}, nil
}
