package app

import (
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/models"
)

type Option func(a *App)

func WithConfig(config env.Config) Option {
	return func(a *App) {
		a.config = config
	}
}

func WithDatabase(db *models.Database) Option {
	return func(a *App) {
		a.db = db
	}
}
