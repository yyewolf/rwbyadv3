package filters

import "github.com/yyewolf/rwbyadv3/internal/cards"

func Or(filters ...SingleFilter) SingleFilter {
	return func(card *cards.Card) bool {
		for _, filter := range filters {
			if filter(card) {
				return true
			}
		}
		return false
	}
}

func And(filters ...SingleFilter) SingleFilter {
	return func(card *cards.Card) bool {
		for _, filter := range filters {
			if !filter(card) {
				return false
			}
		}
		return true
	}
}
