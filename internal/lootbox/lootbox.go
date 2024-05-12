package lootbox

import (
	"math"

	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/internal/filters"
	"github.com/yyewolf/rwbyadv3/internal/stats"
	"github.com/yyewolf/rwbyadv3/models"
)

type Lootbox struct {
	possibilities []*cards.Card

	filter filters.Filter

	ivDie     stats.Dice
	rarityDie stats.Dice
}

func NewLootbox(options ...LootboxOption) *Lootbox {
	l := &Lootbox{}

	for _, option := range options {
		option(l)
	}

	return l
}

func (l *Lootbox) PickCard(m map[string]*cards.Card) *models.Card {
	// Only used the first time around
	if len(l.possibilities) == 0 {
		cards := make([]*cards.Card, 0, len(m))
		for _, card := range m {
			cards = append(cards, card)
		}
		l.possibilities = l.filter(cards)
	}

	n := stats.RandN(0, len(l.possibilities))
	pickedCard := l.possibilities[n]

	// Convert to model
	c := models.Card{
		ID:              uuid.NewString(),
		Level:           1,
		CardType:        string(pickedCard.ID),
		IndividualValue: stats.Roll(l.ivDie),
		Rarity:          int(math.Ceil(stats.Roll(l.rarityDie)/(100.0/6.0)) - 1),
	}

	return &c
}
