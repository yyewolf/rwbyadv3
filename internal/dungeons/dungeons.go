package dungeons

import (
	"math/rand"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/internal/maze"
	"github.com/yyewolf/rwbyadv3/pkg/loottables"
	"github.com/yyewolf/rwbyadv3/pkg/loottables/item"
)

const (
	DefaultHeight = 15
	DefaultWidth  = 15
)

type Dungeon struct {
	r *rand.Rand

	Width  int `json:"width"`
	Height int `json:"height"`

	maze *maze.Grid `json:"-"`
	Grid [][]int    `json:"grid"`

	Loots []loots.DungeonLoot `json:"loots"`
}

func NewDungeon(r *rand.Rand) *Dungeon {
	d := Dungeon{
		r:      r,
		Width:  DefaultWidth,
		Height: DefaultHeight,
	}

	d.maze = maze.NewGrid(r, DefaultHeight, DefaultWidth)
	d.maze.Generate()
	d.Grid = d.maze.Expand(2)

	d.GenerateLoots()

	return &d
}

func (d *Dungeon) GenerateLoots() {
	var possiblePoints [][2]int

	for i := 0; i < len(d.Grid); i++ {
		for j := 0; j < len(d.Grid[i]); j++ {
			if d.Grid[i][j] == 0 {
				possiblePoints = append(possiblePoints, [2]int{i, j})
			}
		}
	}

	var dynamicLootCount = d.Width*d.Height/25 + d.r.Intn(d.Width*d.Height/50)
	var one = float64(dynamicLootCount)

	// Create the loot table
	var lootTable = loottables.New(
		loottables.NewSubLootTable(1,
			loottables.Always(),
			loottables.Unique(),
			loottables.WithCount(dynamicLootCount),
			loottables.WithEntries(
				item.New(&loots.Liens{}, one*4, item.WithAmountRange[int, *loots.Liens](50, 200, 1), item.WithRepartitionFunc[int, *loots.Liens](item.RepartitionGaussian[int](120, 50))),
				// item.New(item.Nothing{}, one),
			),
		),

		loottables.NewSubLootTable(1,
			loottables.Always(),
			loottables.Unique(),
			loottables.WithCount(1),
			loottables.WithEntries(
				item.New(&loots.Exit{}, 1, item.Always[int, *loots.Exit](), item.Unique[int, *loots.Exit]()),
			),
		),
	)

	list := lootTable.ChooseRandomItems(d.r, 2)

	// Place loots
	for _, l := range list {
		if loot, ok := l.(loots.DungeonLoot); ok {
			var location = possiblePoints[d.r.Intn(len(possiblePoints))]
			d.Loots = append(d.Loots, loot.Place(location))
		} else {
			logrus.Errorf("Invalid loot type for dungeons: %T", l)
		}
	}
}
