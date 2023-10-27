package onlinekhabar

import (
	"net/http"
	"strings"

	"nepse-backend/nepse/onlinekhabar"

	"github.com/gin-gonic/gin"
)

type Fundamental struct {
	Ticker          string  `json:"ticker"`
	EPS             float64 `json:"eps"`
	PE              float64 `json:"pe"`
	ROE             float64 `json:"roe"`
	PB              float64 `json:"pb"`
	Beta            float64 `json:"beta"`
	DividendYield   float64 `json:"dividend_yield"`
	Revenue         float64 `json:"revenue"`
	NetProfit       float64 `json:"net_profit"`
	GrossProfit     float64 `json:"gross_profit"`
	NetIncome       float64 `json:"net_income"`
	NetProfitMargin float64 `json:"net_profit_margin"`
	BookValue       float64 `json:"book_value"`
	DebtToEquity    float64 `json:"debt_to_equity"`
	Assets          float64 `json:"assets"`
	Liabilities     float64 `json:"liabilities"`
	Equities        float64 `json:"equities"`
}

// GetFundamental godoc
// @Summary Get Fundamental Details of the Stock
// @Description List all available Nepse Stock
// @Tags nepse
// @Accept  json
// @Produce  json
// @Success 200 {object} []Fundamental
// @Router /leads/info/accounts [get]
func (ok okController) GetFundamental(ctx *gin.Context) {
	sector := ctx.Param("sector")

	sector = stockMap[sector]

	stocks, err := ok.okstock.GetStocks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var tickers []string
	var fundamental []Fundamental
	financialOverviewMap := make(map[string]onlinekhabar.FinancialOverview)
	fundamentalOverviewMap := make(map[string]onlinekhabar.FundamentalOverview)
	balanceSheetMap := make(map[string]onlinekhabar.BalanceSheet)
	incomeStatementMap := make(map[string]onlinekhabar.IncomeStatement)

	for _, stck := range stocks {
		if stck.Sector == sector {
			if !strings.Contains(stck.TickerName, "Promoter") {
				if !strings.Contains(stck.TickerName, "Preferred") {
					tickers = append(tickers, stck.Ticker)
				}
			}
		}
	}

	for _, tick := range tickers {
		fundamental, err := ok.okstock.GetFundamentalQuickView(ctx, tick)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		fundamentalOverviewMap[tick] = *fundamental

		financial, err := ok.okstock.GetFinancialOverview(ctx, tick)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		financialOverviewMap[tick] = *financial

		balance, err := ok.okstock.GetBalanceSheet(ctx, tick)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		balanceSheetMap[tick] = *balance

		income, err := ok.okstock.GetIncomeStatement(ctx, tick)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		incomeStatementMap[tick] = *income
	}

	for _, tick := range tickers {
		revenue := financialOverviewMap[tick].Response.Revenue
		netProfit := financialOverviewMap[tick].Response.Netprofit
		bookValue := financialOverviewMap[tick].Response.Bvps
		debt := financialOverviewMap[tick].Response.Debttoequity
		assets := balanceSheetMap[tick].Response.Totalasset
		liabilities := balanceSheetMap[tick].Response.Totalliabilities
		equity := balanceSheetMap[tick].Response.Totalequities
		grossProfit := incomeStatementMap[tick].Response.Grossprofit
		netIncome := incomeStatementMap[tick].Response.Netincome
		netProfitMargin := incomeStatementMap[tick].Response.Netprofitmargin
		overview := fundamentalOverviewMap[tick]
		latestRevenue := revenue[len(revenue)-1]
		latestNetProfit := netProfit[len(netProfit)-1]
		latestBookValue := bookValue[len(bookValue)-1]
		latestDebt := debt[len(debt)-1]
		latestAssets := assets[len(assets)-1]
		latestLiabilities := liabilities[len(liabilities)-1]
		latestEquity := equity[len(equity)-1]
		latestGrossProfit := grossProfit[len(grossProfit)-1]
		latestNetIncome := netIncome[len(netIncome)-1]
		latestNetProfitMargin := netProfitMargin[len(netProfitMargin)-1]

		fundamental = append(fundamental, Fundamental{
			Ticker:          tick,
			EPS:             overview.Response.EpsDiluted,
			PE:              overview.Response.PeDiluted,
			ROE:             overview.Response.Roe,
			PB:              overview.Response.PbRatio,
			Beta:            overview.Response.Beta,
			DividendYield:   overview.Response.DivYield,
			Revenue:         latestRevenue.Value,
			NetProfit:       latestNetProfit.Value,
			BookValue:       latestBookValue.Value,
			DebtToEquity:    latestDebt.Value,
			Assets:          latestAssets.Value,
			Liabilities:     latestLiabilities.Value,
			Equities:        latestEquity.Value,
			GrossProfit:     latestGrossProfit.Value,
			NetIncome:       latestNetIncome.Value,
			NetProfitMargin: latestNetProfitMargin.Value,
		})
	}

	ctx.JSON(http.StatusOK, fundamental)
}
