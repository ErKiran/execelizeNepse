package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncomeStatement struct {
	Response struct {
		Revenue         []QuarterValue `json:"revenue"`
		Grossprofit     []QuarterValue `json:"grossProfit"`
		Netincome       []QuarterValue `json:"netIncome"`
		Netprofitmargin []QuarterValue `json:"netProfitMargin"`
	} `json:"response"`
}

func (ok *OkStockAPI) GetIncomeStatement(ctx *gin.Context, ticker string) (*IncomeStatement, error) {
	url := ok.buildFundamentalOverview(TickerIncomeStatement, ticker)
	req, err := ok.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res *IncomeStatement
	if _, err := ok.client.Do(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
