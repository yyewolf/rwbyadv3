package utils

import (
	"math"

	"github.com/yyewolf/rwbyadv3/models"
)

type Player struct{}

var Players Player

func (Player) CalculateNextLevelXP(p *models.Player) int64 {
	return int64(10*int(math.Pow(float64(p.Level), 1.8)) + 20)
}

func (c Player) GiveXP(p *models.Player, given int64) (levelUp bool) {
	for p.XP+given >= p.NextLevelXP {
		levelUp = true

		// When leveling up, reiter giving but with lower amount
		given -= p.NextLevelXP - p.XP
		p.Level++
		p.XP = 0
		p.NextLevelXP = c.CalculateNextLevelXP(p)
	}
	p.XP += given
	p.NextLevelXP = c.CalculateNextLevelXP(p)
	return levelUp
}
