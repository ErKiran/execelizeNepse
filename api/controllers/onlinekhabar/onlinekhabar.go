package onlinekhabar

import (
	"nepse-backend/nepse/onlinekhabar"

	"github.com/gin-gonic/gin"
)

var stockMap = map[string]string{
	"bank": "Banking",
}

type okController struct {
	okstock onlinekhabar.OkStock
}

type OKController interface {
	GetStocks(ctx *gin.Context)
	GetFundamental(ctx *gin.Context)
}

func NewOKController() okController {
	return okController{
		okstock: onlinekhabar.NewOkStock(),
	}
}
