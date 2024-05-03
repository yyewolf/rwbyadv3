package interfaces

import (
	"context"

	"github.com/yyewolf/rwbyadv3/models"
)

type Context interface {
	GetPlayer() *models.Player
}

type ContextGenerator interface {
	NewContext(ctx context.Context) Context
}
