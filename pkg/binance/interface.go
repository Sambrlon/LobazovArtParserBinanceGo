package binance

type Binance interface {
	GetTickerPrice(symbol string) (float64, error)
}
