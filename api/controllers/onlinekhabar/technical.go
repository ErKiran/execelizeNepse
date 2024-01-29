package onlinekhabar

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TechnicalCalculation struct {
	MACD int
	RSI  int
	MFI  int
	CCI  int
	MA   int
	ADX  int
}

type Technical struct {
	Ticker string  `json:"ticker"`
	LTP    float64 `json:"ltp"`
	Macd   float64 `json:"macd"`
	Rsi    float64 `json:"rsi"`
	Mfi    float64 `json:"mfi"`
	Cci    float64 `json:"cci"`
	Ma20   float64 `json:"ma20"`
	Ma50   float64 `json:"ma50"`
	Ma200  float64 `json:"ma200"`
	Adx    float64 `json:"adx"`
	Total  float64 `json:"total"`
}

// GetStocks godoc
// @Summary Get Stock List
// @Description List all available Nepse Stock
// @Tags nepse
// @Accept  json
// @Produce  json
// @Success 200 {object} []onlinekhabar.TickerInfo
// @Router /leads/info/accounts [get]
func (ok okController) GetTechnical(ctx *gin.Context) {
	live, err := ok.okstock.GetLiveTrading(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	stocks, err := ok.okstock.GetStocks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var technicals []Technical
	for _, stck := range stocks {
		if stck.Sector == "Mutual Fund" || stck.Sector == "Bond" {
			continue
		}
		if strings.Contains(stck.TickerName, "Promoter") ||
			strings.Contains(stck.TickerName, "Scheme") ||
			strings.Contains(stck.TickerName, "Preferred") || strings.Contains(stck.TickerName, "Lube") {
			continue
		}
		technical, err := ok.okstock.GetTechnicalIndicator(ctx, stck.Ticker)
		if err != nil {
			fmt.Println("unable to get technical indicator for stock", stck.TickerName)
			continue
		}

		if technical == nil {
			continue
		}

		if technical.Response.Macd == 0 || technical.Response.Rsi == 0 {
			continue
		}

		var currentPrice float64

		for _, lv := range live.Response {
			if lv.Ticker == stck.Ticker {
				currentPrice = lv.Ltp
			}
		}

		if currentPrice == 0 {
			continue
		}

		var calcu TechnicalCalculation

		calcu.MACD = ok.CalculateTechnical(technical.Response.Macd, 0, 5, 0, -5)
		calcu.RSI = ok.CalculateTechnical(technical.Response.Rsi, 50, 70, 30, 20)
		calcu.MFI = ok.CalculateTechnical(technical.Response.Mfi, 30, 80, 20, 0)
		calcu.CCI = ok.CalculateTechnical(technical.Response.Cci, 0, 70, 0, -70)
		calcu.ADX = ok.CalculateTechnical(technical.Response.Adx, 15, 25, 10, 5)

		calcu.MA = ok.CalculateMA(currentPrice, technical.Response.Ma20, technical.Response.Ma50, technical.Response.Ma200)

		technicals = append(technicals, Technical{
			Ticker: stck.Ticker,
			LTP:    currentPrice,
			Macd:   technical.Response.Macd,
			Rsi:    technical.Response.Rsi,
			Mfi:    technical.Response.Mfi,
			Cci:    technical.Response.Cci,
			Ma20:   technical.Response.Ma20,
			Ma50:   technical.Response.Ma50,
			Ma200:  technical.Response.Ma200,
			Adx:    technical.Response.Adx,
			Total:  float64(calcu.MACD) + float64(calcu.RSI) + float64(calcu.MFI) + float64(calcu.CCI) + float64(calcu.ADX),
		})
	}

	ctx.JSON(http.StatusOK, technicals)
}

func (ok okController) CalculateMA(currentPrice, ma20, ma50, ma200 float64) int {
	if currentPrice > ma20 && currentPrice > ma50 && currentPrice > ma200 {
		return STRONGBULLISH
	}

	if currentPrice > ma20 || currentPrice > ma50 || currentPrice > ma200 {
		return BULLISH
	}

	if currentPrice < ma20 && currentPrice < ma50 && currentPrice < ma200 {
		return STRONGBEARISH
	}

	if currentPrice < ma20 || currentPrice < ma50 || currentPrice < ma200 {
		return BEARISH
	}
	return 0
}

func (ok okController) CalculateTechnical(indicatorValue float64, upperRange, extremeUpperRange, lowerRange, extremeLowerRange int) int {
	if indicatorValue > float64(upperRange) {
		if indicatorValue > float64(extremeUpperRange) {
			return STRONGBULLISH
		}
		return BULLISH
	}

	if indicatorValue < float64(lowerRange) {
		if indicatorValue < float64(extremeLowerRange) {
			return STRONGBEARISH
		}
		return BEARISH
	}
	return 0
}
