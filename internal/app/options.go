package app

import (
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type Option func(a *App)

func WithConfig(config env.Config) Option {
	return func(a *App) {
		a.config = config
	}
}

func WithDatabase(db interfaces.Database) Option {
	return func(a *App) {
		a.db = db
	}
}
