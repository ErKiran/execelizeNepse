package onlinekhabar

import (
	"nepse-backend/nepse/onlinekhabar"

	"github.com/gin-gonic/gin"
)

var stockMap = map[string]string{
	"bank":    "Banking",
	"hydro":   "HydroPower",
	"finance": "Finance",
	"micro":   "Microfinance",
	"life":    "Life Insurance",
	"trade":   "Trading",
	"manu":    "Manu.& Pro.",
	"invest":  "Investment",
	"non":     "Non Life Insurance",
	"others":  "Others",
	"hotel":   "Hotels And Tourism",
	"dev":     "Development Bank",
}

const (
	NEUTRAL       = 0
	BULLISH       = 1
	STRONGBULLISH = 2
	BEARISH       = -1
	STRONGBEARISH = -2
)

type okController struct {
	okstock onlinekhabar.OkStock
}

type OKController interface {
	GetStocks(ctx *gin.Context)
	GetFundamental(ctx *gin.Context)
	GetTechnical(ctx *gin.Context)
}

func NewOKController() okController {
	return okController{
		okstock: onlinekhabar.NewOkStock(),
	}
}
