package app

import (
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Option func(a *App)

func WithConfig(config *env.Config) Option {
	return func(a *App) {
		a.config = config
	}
}

func WithDatabase(entClient *ent.Client) Option {
	return func(a *App) {
		a.entClient = entClient
	}
}

func WithTemporal(temporal client.Client, w worker.Worker) Option {
	return func(a *App) {
		a.temporalClient = temporal
		a.temporalWorker = w
	}
}

func WithWeb() Option {
	return func(a *App) {
		a.enableWeb = true
	}
}
