package onlinekhabar

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFundamental godoc
// @Summary Get Fundamental Details of the Stock
// @Description List all available Nepse Stock
// @Tags nepse
// @Accept  json
// @Produce  json
// @Success 200 {object} onlinekhabar.FundamentalOverview
// @Router /leads/info/accounts [get]
func (ok okController) GetFundamental(ctx *gin.Context) {
	ticker := ctx.Param("id")
	fundamental, err := ok.okstock.GetFundamentalQuickView(ctx, ticker)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, fundamental)
}
