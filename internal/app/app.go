package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/internal/repo"
	"github.com/yyewolf/rwbyadv3/internal/values"
	"github.com/yyewolf/rwbyadv3/web"

	sloglogrus "github.com/samber/slog-logrus/v2"
)

type App struct {
	config *env.Config

	// discord stuff
	handler *handler.Mux
	client  bot.Client

	ms *builder.MenuStore

	// github stuff
	github *repo.GithubClient

	// jobs stuff
	jobHandler interfaces.JobHandler

	// graceful shutdown
	shutdown     chan struct{}
	errorChannel chan error

	// options
	enableWeb bool
	webApp    *web.WebApp
}

func New(options ...Option) interfaces.App {
	var app = &App{}

	for _, opt := range options {
		opt(app)
	}

	app.jobHandler = jobs.New(
		jobs.WithConfig(app.config),
	)

	app.handler = handler.New()

	c, err := disgo.New(app.config.Discord.Token,
		bot.WithLogger(slog.New(sloglogrus.Option{Level: slog.Level(logrus.GetLevel()), Logger: logrus.StandardLogger()}.NewLogrusHandler())),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds),
		),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagGuilds)),
		bot.WithEventListenerFunc(app.OnReady),
		bot.WithEventListeners(app.handler),
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("could not start discord client")
	}

	app.client = c

	app.shutdown = make(chan struct{})
	app.errorChannel = make(chan error)

	app.github = repo.NewGithubClient(app.config)

	if app.enableWeb {
		app.webApp = web.NewWebApp(
			web.WithApp(app),
		)
	}

	// Jobs
	app.jobHandler.RegisterJobKey("cleanup_db", app.CleanupJob)
	app.jobHandler.ScheduleRecurringJob(
		"cleanup_db",
		time.Date(2024, 1, 1, 1, 1, 0, 0, time.Local),
		24*time.Hour,
	)

	return app
}

func (a *App) OnReady(_ *events.Ready) {
	logrus.Info("Bot is ready, registering commands...")

	a.ms = commands.RegisterCommands(a)

	// Begin job handler here
	go func() {
		err := a.jobHandler.Start()
		if err == nil {
			return
		}

		select {
		case a.errorChannel <- err:
		default:
		}
	}()

	// Set status depending on mode :
	switch a.config.Mode {
	case values.Dev:
		a.client.SetPresence(context.TODO(), gateway.WithPlayingActivity("/help for help - dev"))
	case values.Preprod:
		a.client.SetPresence(context.TODO(), gateway.WithPlayingActivity("/help for help - preprod"))
	case values.Prod:
		a.client.SetPresence(context.TODO(), gateway.WithPlayingActivity("/help for help"))
	}
}

func (a *App) Client() bot.Client {
	return a.client
}

func (a *App) Handler() *handler.Mux {
	return a.handler
}

func (a *App) Config() *env.Config {
	return a.config
}

func (a *App) Github() *repo.GithubClient {
	return a.github
}

func (a *App) JobHandler() interfaces.JobHandler {
	return a.jobHandler
}
