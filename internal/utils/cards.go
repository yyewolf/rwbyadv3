package utils

import (
	"fmt"
	"math"

	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/models"
)

type Card struct{}

var Cards Card

func (Card) Template(c *models.Card) *cards.Card {
	return cards.Cards[c.CardType]
}

func (card Card) GenerateStats(c *models.Card) *models.CardsStat {
	def := card.Template(c)

	return &models.CardsStat{
		CardID:  c.ID,
		Damage:  int(float64(def.BaseStats.Damage) + float64(9*c.Level)*float64(c.IndividualValue/100.0)*math.Pow(2, float64(c.Rarity)/4.6)*math.Pow(3, float64(c.Buffs)/7.0)),
		Healing: int(float64(def.BaseStats.Healing) + float64(11*c.Level)*float64(c.IndividualValue/100.0)*math.Pow(2, float64(c.Rarity)/9.0)*math.Pow(3, float64(c.Buffs)/10.0)),
		Armor:   int(float64(def.BaseStats.Armor) + float64(8*c.Level)*float64(c.IndividualValue/100.0)*math.Pow(2, float64(c.Rarity)/11.8)*math.Pow(3, float64(c.Buffs)/14.0)),
		Health:  int(float64(def.BaseStats.Health) + float64(18*c.Level)*float64(c.IndividualValue/100.0)*math.Pow(2, float64(c.Rarity)/4.6)*math.Pow(3, float64(c.Buffs)/7.0)),
		Speed:   def.BaseStats.Speed,
	}
}

func (card Card) FullString(c *models.Card) string {
	def := card.Template(c)
	return fmt.Sprintf("%s level %d (%d/%dXP) %s (%.2f%%)", card.RarityString(c.Rarity, c.Buffs), c.Level, c.XP, c.NextLevelXP, def.Name, c.IndividualValue)
}

func (card Card) RarityString(rarity, buffs int) (x string) {
	switch rarity {
	case 0: // Common
		x = "□ Common"
	case 1: // Uncommon
		x = "◇ Uncommon"
	case 2: // Rare
		x = "♡ Rare"
	case 3: // Very Rare
		x = "♤ Very Rare"
	case 4: // Legendary
		x = "♧ Legendary"
	case 5: // Collector
		x = "☆ Collector"
	}

	for i := 0; i < buffs; i++ {
		x += "+"
	}
	return x
}
