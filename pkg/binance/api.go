package binance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

type BinanceClient struct {
	client *binance.Client
}

func NewBinanceClient(APIKey, SecretKey string) *BinanceClient {
	return &BinanceClient{
		client: binance.NewClient(APIKey, SecretKey),
	}
}

func (b *BinanceClient) GetTickerPrice(symbol string) (float64, error) {
	// Получение цены используя Binance API
	ctx := context.Background()
	ticker, err := b.client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return 0, err
	}

	if len(ticker) == 0 {
		return 0, fmt.Errorf("ticker not found")
	}

	price, err := strconv.ParseFloat(ticker[0].Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
