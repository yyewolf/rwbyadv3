package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/lib/pq"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
	migrate := dbmate.New(databaseURL)
	migrate.SchemaFile = appConfig.Database.SchemaFile
	migrate.MigrationsDir = []string{appConfig.Database.MigrationsFolder}
	migrate.Log = logrus.New().Writer()

	err := migrate.CreateAndMigrate()
	if err != nil {
		logrus.
			WithError(err).
			Fatal("cannot run migration")
	}

	urlTest := *databaseURL
	q := urlTest.Query()
	q.Add("search_path", "test")
	urlTest.RawQuery = q.Encode()

	entClient, err := ent.Open("postgres", urlTest.String())
	if err != nil {
		logrus.
			WithError(err).
			Fatal("cannot open ent client")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", appConfig.Database.User, appConfig.Database.Pass, appConfig.Database.Database, appConfig.Database.Host, appConfig.Database.Port))
	if err != nil {
		logrus.
			WithError(err).
			Fatal("cannot connect to database")
	}

	boil.SetDB(db)

	cards.ParseCards(appConfig.App.CardsLocation)

	for _, card := range cards.Cards {
		entClient.CardType.Create().
			SetID(card.ID).
			SetName(card.Name).
			SetCategories(card.Categories).
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

	db.Close()
	app.Shutdown()
}
