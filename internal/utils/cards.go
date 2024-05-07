package utils

import (
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
	"github.com/sirupsen/logrus"
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
	return fmt.Sprintf("%s level %d (%d/%dXP) %s (%.2f%%)", card.RarityString(c), c.Level, c.XP, c.NextLevelXP, def.Name, c.IndividualValue)
}

func (card Card) RarityString(c *models.Card) (x string) {
	switch c.Rarity {
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

	for i := 0; i < c.Buffs; i++ {
		x += "+"
	}
	return x
}

func (card Card) RarityToColor(c *models.Card) int {
	EmbedColor := 0
	switch c.Rarity {
	case 0:
		EmbedColor = 0x808080
	case 1:
		EmbedColor = 0x285300
	case 2:
		EmbedColor = 0x00008b
	case 3:
		EmbedColor = 0xB22222
	case 4:
		EmbedColor = 0x800080
	case 5:
		EmbedColor = 0x121212
	}
	return EmbedColor
}

func (card Card) Message(c *models.Card) (*discord.File, discord.Embed, *discord.ContainerComponent) {
	img, err := cards.GetEmbeddableImage(c.CardType, "battle", "png")
	if err != nil {
		logrus.Panic("no image found, safeguard failed :o")
	}

	f := discord.NewFile("ch.png", "", img)

	def := card.Template(c)

	embed := discord.NewEmbedBuilder().
		SetTitlef("Level %d %s", c.Level, def.Name).
		SetColor(card.RarityToColor(c)).
		AddFields(
			discord.EmbedField{
				Name: "**Statistics :**",
				Value: fmt.Sprintf("Category : **%v**\n", def.Categories) +
					fmt.Sprintf("XP : %d/%d\n", c.XP, c.NextLevelXP) +
					fmt.Sprintf("Value : %.2f%%\n", c.IndividualValue) +
					fmt.Sprintf("Rarity : %s\n", card.RarityString(c)) +
					fmt.Sprintf("Health : %d\n", c.R.CardsStat.Health) +
					fmt.Sprintf("Armor : %v\n", c.R.CardsStat.Armor) +
					fmt.Sprintf("Damage : %v\n", c.R.CardsStat.Damage),
			},
		).
		SetThumbnail("attachment://ch.png").
		Build()

	return f, embed, nil
}
