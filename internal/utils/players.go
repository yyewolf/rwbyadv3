package utils

import (
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/env"
)

type Player struct {
	c *env.Config
}

func init() {
	Players = Player{
		c: env.Get(),
	}
}

var Players Player

func (p Player) MaxSlots(player *ent.Player) int64 {
	return player.BackpackLevel * int64(p.c.App.BackpackSize)
}

func (p Player) UsedSlots(player *ent.Player) int64 {
	return int64(len(player.Edges.Cards)) + player.BackpackReservedSlots
}

func (p Player) AvailableSlots(player *ent.Player) int64 {
	return p.MaxSlots(player) - p.UsedSlots(player)
}
