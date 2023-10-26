package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStocks godoc
// @Summary Get Stock List
// @Description List all available Nepse Stock
// @Tags nepse
// @Accept  json
// @Produce  json
// @Success 200 {object} []onlinekhabar.TickerInfo
// @Router /leads/info/accounts [get]
func (ok okController) GetStocks(ctx *gin.Context) {
	stocks, err := ok.okstock.GetStocks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, stocks)
}
