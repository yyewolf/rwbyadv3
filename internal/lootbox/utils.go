package lootbox

import (
	"github.com/yyewolf/rwbyadv3/internal/filters"
	"github.com/yyewolf/rwbyadv3/internal/stats"
)

var (
	NormalLootBox = NewLootbox(
		WithFilter(
			filters.NewFilter(
				filters.All(),
			),
		),
		WithIVDie(stats.NormalLootDie),
		WithRarityDie(stats.NormalLootExp),
	)

	RareLootBox = NewLootbox(
		WithFilter(
			filters.NewFilter(
				filters.All(),
			),
		),
		WithIVDie(stats.RareLootDie),
		WithRarityDie(stats.RareLootExp),
	)

	// TODO: remove this and fix where it was used to implement correct filters
	SpecialLootBox = NewLootbox(
		WithFilter(
			filters.NewFilter(
				filters.All(),
			),
		),
		WithIVDie(stats.LimitedLootDie),
		WithRarityDie(stats.LimitedLootExp),
	)

	// TODO: remove this and fix where it was used to implement correct filters
	LimitedLootBox = NewLootbox(
		WithFilter(
			filters.NewFilter(
				filters.All(),
			),
		),
		WithIVDie(stats.LimitedLootDie),
		WithRarityDie(stats.LimitedLootExp),
	)
)
