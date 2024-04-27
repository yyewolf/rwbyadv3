package app

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/commands"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/models"
	"github.com/yyewolf/rwbyadv3/internal/values"
)

type App struct {
	config env.Config

	// database stuff
	db *models.Database

	// discord stuff
	state  *state.State
	router *cmdroute.Router
	cr     interfaces.CommandRepository

	// graceful shutdown
	shutdown     chan struct{}
	errorChannel chan error
}

func New(options ...Option) interfaces.App {
	var app = &App{}

	for _, opt := range options {
		opt(app)
	}

	app.router = cmdroute.NewRouter()
	app.state = state.New(fmt.Sprintf("Bot %s", app.config.Discord.Token))
	app.state.AddInteractionHandler(app.router)
	app.state.AddIntents(gateway.IntentGuilds)
	app.cr = commands.New(app)

	app.state.AddHandler(app.ReadyHandler)

	app.shutdown = make(chan struct{})
	app.errorChannel = make(chan error)

	return app
}

func (a *App) ReadyHandler(rd *gateway.ReadyEvent) {
	logrus.Info("Bot is ready, registering commands...")
	err := a.cr.RegisterCommands()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to register commands")
	}

	// Set status depending on mode :
	switch a.config.Mode {
	case values.Dev:
		a.state.SendGateway(context.TODO(), &gateway.UpdatePresenceCommand{
			Activities: []discord.Activity{
				{
					Name: "/help for help - dev",
				},
			},
		})
	case values.Preprod:
		a.state.SendGateway(context.TODO(), &gateway.UpdatePresenceCommand{
			Activities: []discord.Activity{
				{
					Name: "/help for help - preprod",
				},
			},
		})
	case values.Prod:
		a.state.SendGateway(context.TODO(), &gateway.UpdatePresenceCommand{
			Activities: []discord.Activity{
				{
					Name: "/help for help",
				},
			},
		})
	}
}

func (a *App) State() *state.State {
	return a.state
}

func (a *App) CommandRouter() *cmdroute.Router {
	return a.router
}

func (a *App) Config() *env.Config {
	return &a.config
}
