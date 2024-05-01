package app

import (
	"github.com/yyewolf/rwbyadv3/internal/env"
)

type Option func(a *App)

func WithConfig(config *env.Config) Option {
	return func(a *App) {
		a.config = config
	}
}

func WithWeb() Option {
	return func(a *App) {
		a.enableWeb = true
	}
}
