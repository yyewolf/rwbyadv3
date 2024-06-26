package utils

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/models"
)

type Card struct{}

var Cards Card

func (Card) Primitive(c *models.Card) *cards.Card {
	return cards.Cards[c.CardType]
}

func (card Card) GenerateStats(c *models.Card) *models.CardsStat {
	def := card.Primitive(c)

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
	def := card.Primitive(c)
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
	case 0: // Common
		EmbedColor = 0x808080
	case 1: // Uncommon
		EmbedColor = 0x7CFC00
	case 2: // Rare
		EmbedColor = 0x87CEEB
	case 3: // Very Rare
		EmbedColor = 0xBA55D3
	case 4: // Legendary
		EmbedColor = 0xFFD700
	case 5: // Collector
		EmbedColor = 0xFF0000
	}
	return EmbedColor
}

func (card Card) Message(c *models.Card) (*discord.File, discord.Embed, *discord.ContainerComponent) {
	img, err := cards.GetEmbeddableImage(c.CardType, "battle", "png")
	if err != nil {
		logrus.Panic("no image found, safeguard failed :o")
	}

	f := discord.NewFile("ch.png", "", img)

	def := card.Primitive(c)

	inline := true

	embed := discord.NewEmbedBuilder().
		SetTitlef("Level %d %s (%d/%d XP)", c.Level, def.Name, c.XP, c.NextLevelXP).
		SetColor(card.RarityToColor(c)).
		AddFields(
			discord.EmbedField{
				Name:   "**General :**",
				Inline: &inline,
				Value: fmt.Sprintf("Category : **%v**\n", strings.Join(def.Categories, ", ")) +
					fmt.Sprintf("Rarity : %s\n", card.RarityString(c)) +
					fmt.Sprintf("Value : %.2f%%\n", c.IndividualValue),
			},
			discord.EmbedField{
				Name:   "**Stats :**",
				Inline: &inline,
				Value: fmt.Sprintf("Health : %d\n", c.R.CardsStat.Health) +
					fmt.Sprintf("Armor : %v\n", c.R.CardsStat.Armor) +
					fmt.Sprintf("Damage : %v\n", c.R.CardsStat.Damage),
			},
		).
		SetThumbnail("attachment://ch.png").
		Build()

	return f, embed, nil
}

func (card Card) IconURI(c *models.Card) string {
	uri, _ := url.JoinPath(Players.c.App.BaseURI, cards.MustGetImageURI(c.CardType, "icon", "webp"))
	return uri
}

type CardMetadata struct {
	Location string
}

func (card Card) SetLocation(c *models.Card, location string) {
	meta := card.GetMeta(c)
	meta.Location = location
	card.SaveMeta(c, meta)
}

func (card Card) GetMeta(c *models.Card) *CardMetadata {
	var meta CardMetadata
	// Load
	c.Metadata.Unmarshal(&meta)
	return &meta
}

func (card Card) SaveMeta(c *models.Card, meta *CardMetadata) {
	c.Metadata.Marshal(&meta)
}

// Leveling
func (card Card) GetNextLevelXP(c *models.Card) int64 {
	return int64(50*c.Level*c.Level + 100)
}

func (card Card) GetXPReward(c *models.Card, multiplier int, boost bool) int64 {
	if c.Level >= 500 {
		return 0
	}
	rint := multiplier * (c.Level)
	add := int64(float64((rand.Intn(26+rint))+15) * (math.Pow(float64(c.Level), 0.72) + 1))
	if boost {
		rint = int(math.Floor(((3.0 / 2.0) * float64(multiplier)) * float64(c.Level)))
		add = int64(float64((rand.Intn(33+rint))+25) * (math.Pow(float64(c.Level), 0.84) + 1))
	}
	return add
}

func (card Card) GiveXP(c *models.Card, XP int64) (levelUp bool) {
	for c.XP+XP > c.NextLevelXP {
		levelUp = true
		//if level up
		XP -= c.NextLevelXP - c.XP
		c.Level++
		c.XP = 0
		c.NextLevelXP = card.GetNextLevelXP(c)
		card.GenerateStats(c)
	}
	c.XP += XP
	c.NextLevelXP = card.GetNextLevelXP(c)
	return levelUp
}
