package loots

import (
	"database/sql"
	"math/rand"

	"github.com/yyewolf/rwbyadv3/models"
)

type Loot interface {
	GetID() string
	GetType() string
	GetX() int
	GetY() int

	Generate(r *rand.Rand, point [2]int) Loot
	PickedUp(tx *sql.Tx, p *models.Player)
	RewardText([]Loot) string
}
