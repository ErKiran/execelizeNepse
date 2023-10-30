package cricket

import (
	"os"

	"nepse-backend/utils"

	"github.com/gin-gonic/gin"
)

const (
	ScoreCard = "pages/match/scorecard?lang=en"
)

type CricInfo interface {
	GetScoreCard(ctx *gin.Context, seriesID, matchID string) (*ScoreCardResponse, error)
}

type CrickInfoAPI struct {
	client *utils.Client
}

func NewCrickInfo() CricInfo {
	client := utils.NewClient(nil, os.Getenv("CRICINFO"), "")

	ci := &CrickInfoAPI{
		client: client,
	}
	return ci
}
