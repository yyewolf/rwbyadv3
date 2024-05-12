package utils

import (
	"github.com/yyewolf/rwbyadv3/models"
)

func (Player) LootBoxesCount(p *models.Player) map[models.LootBoxesType]int {
	var counts = make(map[models.LootBoxesType]int)
	for _, b := range p.R.LootBoxes {
		counts[b.Type]++
	}
	return counts
}

func (Player) TakeFirstLootBoxOf(p *models.Player, t models.LootBoxesType) (*models.LootBox, bool) {
	for _, b := range p.R.LootBoxes {
		if b.Type == t {
			return b, true
		}
	}
	return nil, false
}

func (Player) DeleteBoxFromPlayer(p *models.Player, b *models.LootBox) {
	for i, f := range p.R.LootBoxes {
		if f.ID == b.ID {
			p.R.LootBoxes = append(p.R.LootBoxes[:i], p.R.LootBoxes[i+1:]...)
		}
	}
}
