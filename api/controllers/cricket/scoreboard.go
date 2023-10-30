package cricket

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BattingScorecard struct {
	Player string `json:"player"`
	// BallFaced int    `json:"ball_faced"`
	Runs   int `json:"runs"`
	Fours  int `json:"fours"`
	Sixes  int `json:"sixes"`
	Earned int `json:"earned"`
}

type BowlingScorecard struct {
	Player     string `json:"player"`
	MaidenOver int    `json:"maiden_over"`
	Dot        int    `json:"dot"`
	Wickets    int    `json:"wickets"`
	Balls      int    `json:"balls"`
	// Conceded   int    `json:"conceded"`
	Earned int `json:"earned"`
}

type TotalEarning struct {
	Player string `json:"player"`
	Earned int    `json:"earned"`
}

type ScoreCard struct {
	Batting    []BattingScorecard `json:"batting"`
	Bowling    []BowlingScorecard `json:"bowling"`
	Total      []TotalEarning     `json:"total_earning"`
	GrandTotal int                `json:"grand_total"`
}

func (cc cricinfoController) GetScorecard(ctx *gin.Context) {
	seriesID := ctx.Query("seriesId")
	matchID := ctx.Query("matchId")
	scorecard, err := cc.crickinfo.GetScoreCard(ctx, seriesID, matchID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	var batting []BattingScorecard
	var bowling []BowlingScorecard
	var scorecards ScoreCard
	var totalEarn []TotalEarning
	var grandTotal int
	total := make(map[string]int)
	for _, ing := range scorecard.Content.Innings {
		if ing.Team.Name == "Nepal" && ing.Isbatted {
			for _, bat := range ing.Inningbatsmen {
				if bat.Battedtype == "yes" {
					batting = append(batting, BattingScorecard{
						Player: bat.Player.Longname,
						// BallFaced: int(bat.Balls),
						Runs:   int(bat.Runs),
						Fours:  int(bat.Fours),
						Sixes:  int(bat.Sixes),
						Earned: cc.CalculateBatsmenEarning(int(bat.Runs), int(bat.Fours), int(bat.Sixes)),
					})
				}
			}
		}

		if ing.Team.Name != "Nepal" {
			for _, bowl := range ing.Inningbowlers {
				if bowl.Bowledtype == "yes" {
					bowling = append(bowling, BowlingScorecard{
						Player:     bowl.Player.Longname,
						Dot:        int(bowl.Dots),
						MaidenOver: int(bowl.Maidens),
						Wickets:    int(bowl.Wickets),
						Balls:      int(bowl.Balls),
						// Conceded:   int(bowl.Conceded),
						Earned: cc.CalculateBowlerEarning(int(bowl.Dots), int(bowl.Maidens), int(bowl.Wickets)),
					})
				}
			}
		}
	}

	for _, earn := range batting {
		total[earn.Player] += earn.Earned
	}

	for _, earn := range bowling {
		total[earn.Player] += earn.Earned
	}

	for k, v := range total {
		grandTotal += v
		totalEarn = append(totalEarn, TotalEarning{
			Player: k,
			Earned: v,
		})
	}

	scorecards.Batting = batting
	scorecards.Bowling = bowling
	scorecards.Total = totalEarn
	scorecards.GrandTotal = grandTotal
	ctx.JSON(http.StatusOK, scorecards)
}

func (cc cricinfoController) CalculateBatsmenEarning(runs, fours, sixes int) int {
	return runs*RUN + fours*FOUR + sixes*SIX
}

func (cc cricinfoController) CalculateBowlerEarning(dot, maiden, wicket int) int {
	return dot*DOT + maiden*MAIDEN + wicket*WICKET
}
