package tokens

import (
	"time"
)

type Repository interface {
	AddTicker(symbol string, price float64) error
	AddTickerPrice(symbol string, price float64) error
	GetTickerData(symbol string, dateFrom, dateTo time.Time) ([]FetchTickerResponse, error)
	GetPrivateKeyByPublicKey(publicKey string) (string, error)
}
