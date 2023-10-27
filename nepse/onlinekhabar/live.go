package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LiveTradingResponse struct {
	Response []struct {
		Ticker           string  `json:"ticker"`
		Indices          string  `json:"indices"`
		TickerName       string  `json:"ticker_name"`
		Ltp              float64 `json:"ltp"`
		Ltv              float64 `json:"ltv"`
		PointChange      float64 `json:"point_change"`
		PercentageChange float64 `json:"percentage_change"`
		Open             float64 `json:"open"`
		High             float64 `json:"high"`
		Low              float64 `json:"low"`
		Volume           float64 `json:"volume"`
		Previousclosing  float64 `json:"previousClosing"`
		CalculatedOn     string  `json:"calculated_on"`
		Amount           float64 `json:"amount"`
		Datasource       string  `json:"datasource"`
		Icon             string  `json:"icon"`
	} `json:"response"`
}

func (ok *OkStockAPI) GetLiveTrading(ctx *gin.Context) (*LiveTradingResponse, error) {
	req, err := ok.client.NewRequest(http.MethodGet, LiveTrading, nil)
	if err != nil {
		return nil, err
	}

	var res *LiveTradingResponse
	if _, err := ok.client.Do(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
