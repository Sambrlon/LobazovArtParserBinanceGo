package tokens

import "time"

type UseCase interface {
	AddTicker(symbol string) error
	RunFetch(symbol string)
	FetchTickerData(symbol string, dateFrom, dateTo time.Time) ([]FetchTickerResponse, error)
}
