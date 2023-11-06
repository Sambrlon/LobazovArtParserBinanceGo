package tokens

import "time"

type AddTickerParams struct {
	Ticker string `json:"ticker"`
}

type FetchTickerParams struct {
	Ticker   string    `query:"ticker"`
	DateFrom time.Time `query:"date_from"`
	DateTo   time.Time `query:"date_to"`
}

type FetchTickerResponse struct {
	Symbol     string    `db:"symbol"`
	DateTime   time.Time `db:"date_time"`
	Price      float64   `db:"price"`
	Difference float64   `db:"difference"`
}

type APIKey struct {
	ID          uint   `gorm:"primaryKey"`
	PublicKey   string `gorm:"unique;not null"`
	PrivateKey  string `gorm:"not null"`
	Description string
}
