package loots

import (
	"database/sql"

	"github.com/yyewolf/rwbyadv3/ent"
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
	NewPickedUp(tx *ent.Tx, p *ent.Player) error
}

var DungeonLoots = []DungeonLoot{
	Liens{},
}
