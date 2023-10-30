package onlinekhabar

import (
	"os"

	"nepse-backend/utils"

	"github.com/gin-gonic/gin"
)

const (
	ListTicker                = "search-list/tickers"
	TickerFundamentalOverview = "ticker-page/ticker-quick-view"
	TickerFinancialOverview   = "ticker-page/financial-overview"
	TickerBalanceSheet        = "ticker-page/financial-balance-sheet"
	TickerIncomeStatement     = "ticker-page/financial-income-statement"
	LiveTrading               = "stock_live/live-trading"
	Technical                 = "ticker-page/ticker-technical-indicator"
)

type OkStock interface {
	GetStocks(ctx *gin.Context) ([]TickerInfo, error)
	GetFundamentalQuickView(ctx *gin.Context, ticker string) (*FundamentalOverview, error)
	GetFinancialOverview(ctx *gin.Context, ticker string) (*FinancialOverview, error)
	GetBalanceSheet(ctx *gin.Context, ticker string) (*BalanceSheet, error)
	GetIncomeStatement(ctx *gin.Context, ticker string) (*IncomeStatement, error)
	GetLiveTrading(ctx *gin.Context) (*LiveTradingResponse, error)
	GetTechnicalIndicator(ctx *gin.Context, ticker string) (*TechnicalIndicator, error)
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
