package middleware

import (
	"context"

	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

type ContextImpl struct {
	Player *models.Player

	interfaces.Database
	context.Context
}

func (c *ContextImpl) GetPlayer() *models.Player {
	return c.Player
}

type ContextGeneratorImpl struct {
	db interfaces.Database
	context.Context
}

func NewContextGenerator(app interfaces.App) interfaces.ContextGenerator {
	return &ContextGeneratorImpl{
		db: app.Database(),
	}
}

func (cg *ContextGeneratorImpl) NewContext(ctx context.Context) interfaces.Context {
	return &ContextImpl{
		Context: ctx,
	}
}

func ContextualizeCommandRouter[K, T any](cg interfaces.ContextGenerator, f func(ctx interfaces.Context, data K) T) func(ctx context.Context, data K) T {
	return func(ctx context.Context, data K) T {
		return f(cg.NewContext(ctx), data)
	}
}

func ContextualizeEventHandler[K any](cg interfaces.ContextGenerator, f func(ctx interfaces.Context, event K)) func(event K) {
	return func(event K) {
		f(cg.NewContext(context.Background()), event)
	}
}
