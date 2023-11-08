package onlinekhabar

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"nepse-backend/nepse/onlinekhabar"

	"github.com/gin-gonic/gin"
)

func (ok okController) GetTechnicalCon(ctx *gin.Context) {
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
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, stck := range stocks {
		if shouldSkipStock(stck) {
			continue
		}

		wg.Add(1)
		go func(stck onlinekhabar.TickerInfo) {
			defer wg.Done()
			technical, err := ok.okstock.GetTechnicalIndicator(ctx, stck.Ticker)
			if err != nil {
				mu.Lock()
				fmt.Println("Unable to get technical indicator for stock", stck.TickerName)
				mu.Unlock()
				return
			}

			if technical == nil || technical.Response.Macd == 0 || technical.Response.Rsi == 0 {
				return
			}

			currentPrice := getCurrentPrice(live, stck.Ticker)
			if currentPrice == 0 {
				return
			}

			calcu := calculateTechnicalValues(ok, technical, currentPrice)

			mu.Lock()
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
			mu.Unlock()
		}(stck)
	}

	wg.Wait()

	ctx.JSON(http.StatusOK, technicals)
}

func shouldSkipStock(stck onlinekhabar.TickerInfo) bool {
	return stck.Sector == "Mutual Fund" || stck.Sector == "Bond" ||
		strings.Contains(stck.TickerName, "Promoter") ||
		strings.Contains(stck.TickerName, "Scheme") ||
		strings.Contains(stck.TickerName, "Preferred") || strings.Contains(stck.TickerName, "Lube")
}

func getCurrentPrice(liveTradingResponse *onlinekhabar.LiveTradingResponse, ticker string) float64 {
	for _, lv := range liveTradingResponse.Response {
		if lv.Ticker == ticker {
			return lv.Ltp
		}
	}
	return 0
}

func calculateTechnicalValues(ok okController, res *onlinekhabar.TechnicalIndicator, currentPrice float64) TechnicalCalculation {
	response := res.Response
	calcu := TechnicalCalculation{
		MACD: ok.CalculateTechnical(response.Macd, 0, 5, 0, -5),
		RSI:  ok.CalculateTechnical(response.Rsi, 50, 70, 30, 20),
		MFI:  ok.CalculateTechnical(response.Mfi, 30, 80, 20, 0),
		CCI:  ok.CalculateTechnical(response.Cci, 0, 70, 0, -70),
		ADX:  ok.CalculateTechnical(response.Adx, 15, 25, 10, 5),
		MA:   ok.CalculateMA(currentPrice, response.Ma20, response.Ma50, response.Ma200),
	}
	return calcu
}
