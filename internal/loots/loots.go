package loots

import (
	"database/sql"

	"github.com/yyewolf/rwbyadv3/models"
)

type DungeonLoot interface {
	GetID() string
	GetType() string
	GetX() int
	GetY() int

	Place(point [2]int) DungeonLoot

	Loot
}

type Loot interface {
	RewardText([]interface{}) string
	PickedUp(tx *sql.Tx, p *models.Player)
}

var DungeonLoots = []DungeonLoot{
	Liens{},
}
