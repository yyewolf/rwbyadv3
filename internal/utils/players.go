package utils

import (
	"math"
	"math/rand"

	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/models"
)

type Player struct {
	c *env.Config
}

func init() {
	Players = Player{
		c: env.Get(),
	}
}

var Players Player

// Replace with iterator with go1.23
func (Player) AvailableCards(p *models.Player) []*models.Card {
	var cards []*models.Card
	for _, pc := range p.R.PlayerCards {
		if !pc.R.Card.Available {
			continue
		}
		cards = append(cards, pc.R.Card)
	}
	return cards
}

func (Player) GetAvailableCard(p *models.Player, i int) (*models.Card, bool) {
	var at int
	for _, pc := range p.R.PlayerCards {
		if !pc.R.Card.Available {
			continue
		}
		if at == i {
			return pc.R.Card, true
		}
		at++
	}
	return nil, false
}

func (Player) MarketListings(p *models.Player) []*models.Card {
	var cards []*models.Card
	for _, pc := range p.R.PlayerCards {
		if pc.R.Card.Available {
			continue
		}

		// Get metadata
		meta := Cards.GetMeta(pc.R.Card)
		if meta.Location != "listings" {
			continue
		}

		cards = append(cards, pc.R.Card)
	}
	return cards
}

func (Player) GetMarketListing(p *models.Player, i int) (*models.Card, bool) {
	var at int
	for _, pc := range p.R.PlayerCards {
		if !pc.R.Card.Available {
			continue
		}

		// Get metadata
		meta := Cards.GetMeta(pc.R.Card)
		if meta.Location != "listings" {
			continue
		}

		if at == i {
			return pc.R.Card, true
		}

		at++
	}
	return nil, false
}

func (Player) MarketAuctions(p *models.Player) []*models.Card {
	var cards []*models.Card
	for _, pc := range p.R.PlayerCards {
		if pc.R.Card.Available {
			continue
		}

		// Get metadata
		meta := Cards.GetMeta(pc.R.Card)
		if meta.Location != "auctions" {
			continue
		}

		cards = append(cards, pc.R.Card)
	}
	return cards
}

func (Player) GetMarketAuction(p *models.Player, i int) (*models.Card, bool) {
	var at int
	for _, pc := range p.R.PlayerCards {
		if !pc.R.Card.Available {
			continue
		}

		// Get metadata
		meta := Cards.GetMeta(pc.R.Card)
		if meta.Location != "auctions" {
			continue
		}

		if at == i {
			return pc.R.Card, true
		}

		at++
	}
	return nil, false
}

func (Player) AvailableBalance(p *models.Player) int64 {
	return p.Liens - p.LiensBidded
}

func (p Player) MaxSlots(player *models.Player) int {
	return player.BackpackLevel * p.c.App.BackpackSize
}

func (p Player) UsedSlots(player *models.Player) int {
	return len(player.R.PlayerCards) + player.SlotsReserved
}

func (p Player) AvailableSlots(player *models.Player) int {
	return p.MaxSlots(player) - p.UsedSlots(player)
}

// Leveling
func (Player) GetNextLevelXP(p *models.Player) int64 {
	return int64(10*int(math.Pow(float64(p.Level), 1.8)) + 20)
}

func (Player) GetXPReward(p *models.Player, difficulty float64, boost bool) int64 {
	rint := int(5*difficulty*math.Pow(float64(p.Level), 1.48)) + 10
	add := difficulty*float64(rand.Intn(rint)) + 5 + math.Pow(float64(p.Level), 1.45)
	if boost {
		rint = int(((3 / 2) * difficulty) * float64(p.Level))
		add = float64((rand.Intn(33+rint))+25) * (math.Pow(float64(p.Level), 0.84) + 1)
	}
	return int64(add)
}

func (pl Player) GiveXP(p *models.Player, XP int64) (levelUp bool) {
	for p.XP+XP > p.NextLevelXP {
		levelUp = true
		//if level up
		XP -= p.NextLevelXP - p.XP
		p.Level++
		p.XP = 0
		p.NextLevelXP = pl.GetNextLevelXP(p)
	}
	p.XP += XP
	p.NextLevelXP = pl.GetNextLevelXP(p)
	return levelUp
}
