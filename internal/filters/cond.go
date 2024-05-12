package filters

import (
	"slices"

	"github.com/yyewolf/rwbyadv3/internal/cards"
)

func All() SingleFilter {
	return func(card *cards.Card) bool {
		return true
	}
}

func HasCategory(category string) SingleFilter {
	return func(card *cards.Card) bool {
		return slices.Contains(card.Categories, category)
	}
}

func OnlyCategoryIs(category string) SingleFilter {
	return func(card *cards.Card) bool {
		return len(card.Categories) == 1 && card.Categories[0] == category
	}
}
