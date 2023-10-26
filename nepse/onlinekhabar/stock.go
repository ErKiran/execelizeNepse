package onlinekhabar

import (
	"context"
	"net/http"
)

type StockList struct {
	Response []ListStockResponse `json:"response"`
}

type ListStockResponse struct {
	Ticker           string  `json:"ticker"`
	TickerName       string  `json:"ticker_name"`
	Icon             string  `json:"icon"`
	Sector           string  `json:"sector"`
	PointChange      float64 `json:"point_change"`
	PercentageChange float64 `json:"percentage_change"`
}

type TickerInfo struct {
	Ticker     string `json:"ticker"`
	TickerName string `json:"ticker_name"`
	Sector     string `json:"sector"`
}

func (ok *OkStockAPI) GetStocks() ([]TickerInfo, error) {
	req, err := ok.client.NewRequest(http.MethodGet, ListTicker, nil)
	if err != nil {
		return nil, err
	}

	res := &StockList{}
	if _, err := ok.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}

	var stocks []TickerInfo

	for _, ticker := range res.Response {
		stocks = append(stocks, TickerInfo{
			Ticker:     ticker.Ticker,
			TickerName: ticker.TickerName,
			Sector:     ticker.Sector,
		})
	}
	return stocks, nil
}
