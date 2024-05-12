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
)
