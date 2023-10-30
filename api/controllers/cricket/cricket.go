package cricket

import (
	"nepse-backend/nepse/cricket"

	"github.com/gin-gonic/gin"
)

type cricinfoController struct {
	crickinfo cricket.CricInfo
}

const (
	FOUR   = 10000
	SIX    = 20000
	DOT    = 1000
	RUN    = 1000
	WICKET = 25000
	MAIDEN = 10000
)

type CrickInfoController interface {
	GetScorecard(ctx *gin.Context)
}

func NewCricInfoController() cricinfoController {
	return cricinfoController{
		crickinfo: cricket.NewCrickInfo(),
	}
}
