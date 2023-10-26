package onlinekhabar

import (
	"os"

	"nepse-backend/utils"

	"github.com/gin-gonic/gin"
)

const (
	ListTicker                = "search-list/tickers"
	TickerFundamentalOverview = "ticker-page/ticker-quick-view"
)

type OkStock interface {
	GetStocks(ctx *gin.Context) ([]TickerInfo, error)
	GetFundamentalQuickView(ctx *gin.Context, ticker string) (*FundamentalOverview, error)
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
