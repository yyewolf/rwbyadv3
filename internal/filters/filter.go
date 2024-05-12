package filters

import "github.com/yyewolf/rwbyadv3/internal/cards"

type SingleFilter func(*cards.Card) bool
type Filter func([]*cards.Card) []*cards.Card

func NewFilter(filters ...SingleFilter) Filter {
	return func(c []*cards.Card) []*cards.Card {
		var result []*cards.Card
		for _, card := range c {
			if And(filters...)(card) {
				result = append(result, card)
			}
		}
		return result
	}
}
