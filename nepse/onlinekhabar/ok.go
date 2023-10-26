package onlinekhabar

import (
	"os"

	"nepse-backend/utils"
)

const (
	ListTicker = "search-list/tickers"
)

type OkStock interface {
	GetStocks() ([]TickerInfo, error)
}

type OkStockAPI struct {
	client *utils.Client
}

func NewOkStock() OkStock {
	client := utils.NewClient(nil, os.Getenv("ONLINEKHABAR"), "")

	ok := &OkStockAPI{
		client: client,
	}
	return ok
}
