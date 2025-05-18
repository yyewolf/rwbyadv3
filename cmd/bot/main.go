package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/lib/pq"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/cardtype"
	"github.com/yyewolf/rwbyadv3/internal/app"
	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/hooks"
)

func main() {
	env.Load()
	appConfig := env.Get()

	databaseURL, _ := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", appConfig.Database.User, appConfig.Database.Pass, appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.Database))

	entClient, err := ent.Open("postgres", databaseURL.String())
	if err != nil {
		logrus.
			WithError(err).
			Fatal("cannot open ent client")
	}

	cards.ParseCards(appConfig.App.CardsLocation)

	for _, card := range cards.Cards {
		entClient.CardType.Create().
			SetID(card.ID).
			SetName(card.Name).
			SetCategories(card.Categories).
			SetSCategories(strings.Join(card.Categories, ",")).
			OnConflictColumns(cardtype.FieldID).
			UpdateNewValues().
			Exec(context.Background())
	}

	// Create the temporal client
	temporal, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", appConfig.Temporal.Host, appConfig.Temporal.Port),
		Logger:   slog.New(sloglogrus.Option{Logger: logrus.StandardLogger()}.NewLogrusHandler()),
	})
	if err != nil {
		logrus.
			WithError(err).
			Fatal("cannot connect to temporal")
	}

	temporalWorker := worker.New(temporal, appConfig.Temporal.TaskQueue, worker.Options{})

	app := app.New(
		app.WithConfig(appConfig),
		app.WithWeb(),
		app.WithTemporal(temporal, temporalWorker),
		app.WithDatabase(entClient),
	)

	hooks.RegisterHooks(app)

	go app.Start()

	// Listen for CTRL+C
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	logrus.Info("Bot is now running. Press CTRL+C to exit.")
	<-done // Will block here until user hits ctrl+c

	entClient.Close()
	app.Shutdown()
}
