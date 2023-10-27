package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuarterValue struct {
	YearQuarter string  `json:"year_quarter"`
	Value       float64 `json:"value"`
}

type FinancialOverview struct {
	Response struct {
		Revenue      []QuarterValue `json:"revenue"`
		Netprofit    []QuarterValue `json:"netProfit"`
		Eps          []QuarterValue `json:"eps"`
		Bvps         []QuarterValue `json:"bvps"`
		Roe          []QuarterValue `json:"roe"`
		Debttoequity []QuarterValue `json:"debtToEquity"`
	} `json:"response"`
}

func (ok *OkStockAPI) GetFinancialOverview(ctx *gin.Context, ticker string) (*FinancialOverview, error) {
	url := ok.buildFundamentalOverview(TickerFinancialOverview, ticker)
	req, err := ok.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res *FinancialOverview
	if _, err := ok.client.Do(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
