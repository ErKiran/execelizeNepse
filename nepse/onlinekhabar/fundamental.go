package onlinekhabar

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FundamentalOverviewResponse struct {
	Ticker        string  `json:"ticker"`
	Open          float64 `json:"open"`
	EpsDiluted    float64 `json:"eps_diluted"`
	PeDiluted     float64 `json:"pe_diluted"`
	Roe           float64 `json:"roe"`
	PbRatio       float64 `json:"pb_ratio"`
	Beta          float64 `json:"beta"`
	Sector        string  `json:"Sector"`
	DivYield      float64 `json:"div_yield"`
	AvgVolume7Day float64 `json:"avg_volume_7_day"`
}

type FundamentalOverview struct {
	Response FundamentalOverviewResponse `json:"response"`
}

func (ok *OkStockAPI) GetFundamentalQuickView(ctx *gin.Context, ticker string) (*FundamentalOverview, error) {
	url := ok.buildFundamentalOverview(TickerFundamentalOverview, ticker)
	req, err := ok.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var res *FundamentalOverview
	if _, err := ok.client.Do(ctx, req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (ok *OkStockAPI) buildFundamentalOverview(urlPath, ticker string) string {
	return fmt.Sprintf("%s/%s", urlPath, ticker)
}
