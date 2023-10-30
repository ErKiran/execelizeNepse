package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TechnicalIndicator struct {
	Response struct {
		Ticker string  `json:"ticker"`
		Macd   float64 `json:"macd"`
		Rsi    float64 `json:"rsi"`
		Mfi    float64 `json:"mfi"`
		Cci    float64 `json:"cci"`
		Ma20   float64 `json:"ma20"`
		Ma50   float64 `json:"ma50"`
		Ma200  float64 `json:"ma200"`
		Adx    float64 `json:"adx"`
	} `json:"response"`
}

func (ok *OkStockAPI) GetTechnicalIndicator(ctx *gin.Context, ticker string) (*TechnicalIndicator, error) {
	url := ok.buildFundamentalOverview(Technical, ticker)
	req, err := ok.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res *TechnicalIndicator
	if _, err := ok.client.Do(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
