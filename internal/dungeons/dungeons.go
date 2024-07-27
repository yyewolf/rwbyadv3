package dungeons

import (
	"math/rand"

	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/internal/maze"
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

	Loots []loots.Loot `json:"loots"`
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

	// Decide how many loots to generate
	var lootCount = d.Width*d.Height/25 + d.r.Intn(d.Width*d.Height/50)

	for i := 0; i < lootCount; i++ {
		if len(possiblePoints) == 0 {
			break
		}

		var idx = d.r.Intn(len(possiblePoints))
		var lootLocation = possiblePoints[idx]

		// Pick a random loot
		var loot = loots.DungeonLoots[d.r.Intn(len(loots.DungeonLoots))]
		loot = loot.Generate(d.r, lootLocation)

		d.Loots = append(d.Loots, loot)
		possiblePoints = append(possiblePoints[:idx], possiblePoints[idx+1:]...)
	}

	l := loots.DungeonLoots[0].Generate(d.r, [2]int{3, 3})
	d.Loots = append(d.Loots, l)

	// Add exit to loots
	var location = possiblePoints[d.r.Intn(len(possiblePoints))]
	var exit = loots.Exit{}.Generate(d.r, location)

	// Append the exit randomly
	i := d.r.Intn(len(d.Loots))
	d.Loots = append(d.Loots[:i], append([]loots.Loot{exit}, d.Loots[i:]...)...)
}
