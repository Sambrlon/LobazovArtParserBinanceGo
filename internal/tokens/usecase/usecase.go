package usecase

import (
	"errors"
	"log"
	"sambrlon/internal/tokens"
	"sambrlon/pkg/binance"
	"time"
)

type Usecase struct {
	DB      tokens.Repository
	Binance binance.Binance
}

var ErrFetchTicker = errors.New("fetching ticker price error")
var ErrSaveTicker = errors.New("saving ticker price error")

func NewUsecase(db tokens.Repository, binanceClient binance.Binance) *Usecase {
	return &Usecase{
		DB:      db,
		Binance: binanceClient,
	}
}

func (u *Usecase) RunFetch(symbol string) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			log.Println("Fetching ticker price for symbol", symbol)
			price, err := u.Binance.GetTickerPrice(symbol)
			if err != nil {
				log.Println("Error fetching ticker price", err)
				continue
			}

			log.Println("Fetched successfully for symbol", symbol)

			if err := u.DB.AddTickerPrice(symbol, price); err != nil {
				log.Println("Error saving ticker price", err)
			}
		}

	}()
}

//func (u *Usecase) RunFetch(symbol string) {
//	ticker := time.NewTicker(time.Minute)
//	defer ticker.Stop()
//
//	for {
//		<-ticker.C
//		log.Println("Fetching ticker price for symbol:", symbol)
//		price, err := u.Binance.GetTickerPrice(symbol)
//		if err != nil {
//			log.Println("Error fetching ticker price:", err)
//			continue
//		}
//		log.Println("Fetched successfully for symbol:", symbol)
//
//		if err := u.DB.AddTickerPrice(symbol, price); err != nil {
//			log.Println("Error saving ticker price:", err)
//		}
//	}
//}

func (u *Usecase) AddTicker(symbol string) error {
	log.Println("Fetching ticker price...")
	price, err := u.Binance.GetTickerPrice(symbol)
	if err != nil {
		log.Println("Error fetching ticker price:", err)
		return ErrFetchTicker
	}
	log.Println("Fetched successfully")

	if err := u.DB.AddTicker(symbol, price); err != nil {
		log.Println("Error saving ticker price:", err)
		return ErrSaveTicker
	}

	go u.RunFetch(symbol)

	return nil
}

func (u *Usecase) FetchTickerData(symbol string, dateFrom, dateTo time.Time) ([]tokens.FetchTickerResponse, error) {
	data, err := u.DB.GetTickerData(symbol, dateFrom, dateTo)
	if err != nil {
		log.Println("Error fetching ticker data:", err)
		return nil, ErrFetchTicker
	}
	log.Println(data)

	var prevPrice float64
	for i := range data {
		if prevPrice != 0 {
			data[i].Difference = (data[i].Price - prevPrice) / prevPrice * 100
		} else {
			data[i].Difference = 0
		}
		prevPrice = data[i].Price
	}

	return data, nil
}
