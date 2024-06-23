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
