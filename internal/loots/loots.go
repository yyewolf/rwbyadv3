package loots

import (
	"github.com/yyewolf/rwbyadv3/ent"
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
	PickedUp(tx *ent.Tx, p *ent.Player) error
}

var DungeonLoots = []DungeonLoot{
	Liens{},
}
