package interfaces

import (
	"context"

	"github.com/yyewolf/rwbyadv3/ent"
)

type Context interface {
	GetPlayer() *ent.Player
}

type ContextGenerator interface {
	NewContext(ctx context.Context) Context
}
