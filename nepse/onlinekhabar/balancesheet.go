package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BalanceSheet struct {
	Response struct {
		Totalasset       []QuarterValue `json:"totalAsset"`
		Totalliabilities []QuarterValue `json:"totalLiabilities"`
		Totalequities    []QuarterValue `json:"totalEquities"`
	} `json:"response"`
}

func (ok *OkStockAPI) GetBalanceSheet(ctx *gin.Context, ticker string) (*BalanceSheet, error) {
	url := ok.buildFundamentalOverview(TickerBalanceSheet, ticker)
	req, err := ok.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res *BalanceSheet
	if _, err := ok.client.Do(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
