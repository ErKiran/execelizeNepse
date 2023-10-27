package onlinekhabar

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"nepse-backend/nepse/onlinekhabar"

	"github.com/gin-gonic/gin"
)

type Fundamental struct {
	Ticker          string  `json:"ticker"`
	LTP             float64 `json:"ltp"`
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

	liveTrading, err := ok.okstock.GetLiveTrading(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
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
		var ltp float64

		for _, lt := range liveTrading.Response {
			if lt.Ticker == tick {
				ltp = lt.Ltp
			}
		}
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
		if len(revenue) > 0 &&
			len(netProfit) > 0 &&
			len(bookValue) > 0 &&
			len(debt) > 0 &&
			len(assets) > 0 &&
			len(liabilities) > 0 &&
			len(equity) > 0 &&
			len(grossProfit) > 0 &&
			len(netIncome) > 0 &&
			len(netProfitMargin) > 0 {
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
				LTP:             ltp,
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
	}

	if err := ok.FundamentalToCSV(fundamental, sector); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, fundamental)
}

func (ok *okController) FundamentalToCSV(data []Fundamental, sector string) error {
	// Create a CSV file
	file, err := os.Create(fmt.Sprintf("fundamental-%s", sector))
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header
	header := []string{
		"Ticker", "LTP", "EPS", "PE", "ROE", "PB", "Beta",
		"Dividend Yield", "Revenue", "Net Profit", "Gross Profit",
		"Net Income", "Net Profit Margin", "Book Value", "Debt to Equity",
		"Assets", "Liabilities", "Equities",
	}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return err
	}

	// Write the stock data to the CSV file

	for _, stock := range data {
		record := []string{
			stock.Ticker,
			fmt.Sprintf("%.4f", stock.LTP),
			fmt.Sprintf("%.4f", stock.EPS),
			fmt.Sprintf("%.4f", stock.PE),
			fmt.Sprintf("%.4f", stock.ROE),
			fmt.Sprintf("%.4f", stock.PB),
			fmt.Sprintf("%.4f", stock.Beta),
			fmt.Sprintf("%.10f", stock.DividendYield),
			fmt.Sprintf("%.2f", stock.Revenue),
			fmt.Sprintf("%.2f", stock.NetProfit),
			fmt.Sprintf("%.2f", stock.GrossProfit),
			fmt.Sprintf("%.2f", stock.NetIncome),
			fmt.Sprintf("%.4f", stock.NetProfitMargin),
			fmt.Sprintf("%.4f", stock.BookValue),
			fmt.Sprintf("%.10f", stock.DebtToEquity),
			fmt.Sprintf("%.1f", stock.Assets),
			fmt.Sprintf("%.1f", stock.Liabilities),
			fmt.Sprintf("%.1f", stock.Equities),
		}
		err = writer.Write(record)
		if err != nil {
			fmt.Println("Error writing CSV record:", err)
			return err
		}
	}
	return nil
}
