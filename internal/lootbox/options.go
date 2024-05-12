package lootbox

import (
	"github.com/yyewolf/rwbyadv3/internal/filters"
	"github.com/yyewolf/rwbyadv3/internal/stats"
)

type LootboxOption func(*Lootbox)

func WithFilter(filter filters.Filter) func(*Lootbox) {
	return func(l *Lootbox) {
		l.filter = filter
	}
}

func WithIVDie(die stats.Dice) func(*Lootbox) {
	return func(l *Lootbox) {
		l.ivDie = die
	}
}

func WithRarityDie(die stats.Dice) func(*Lootbox) {
	return func(l *Lootbox) {
		l.rarityDie = die
	}
}
