package main

import (
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
	"github.com/yyewolf/rwbyadv3/internal/app"
	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/hooks"
)

func main() {
	env.Load()
	c := env.Get()

	u, _ := url.Parse(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Database.User, c.Database.Pass, c.Database.Host, c.Database.Port, c.Database.Database))
	migrate := dbmate.New(u)
	migrate.SchemaFile = c.Database.SchemaFile
	migrate.MigrationsDir = []string{c.Database.MigrationsFolder}
	migrate.Log = logrus.New().Writer()

	err := migrate.CreateAndMigrate()
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", c.Database.User, c.Database.Pass, c.Database.Database, c.Database.Host, c.Database.Port))
	if err != nil {
		logrus.Fatal(err)
	}

	boil.SetDB(db)

	cards.ParseCards(c.App.CardsLocation)

	// Create the temporal client
	temporal, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", c.Temporal.Host, c.Temporal.Port),
		Logger:   slog.New(sloglogrus.Option{Logger: logrus.StandardLogger()}.NewLogrusHandler()),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	w := worker.New(temporal, c.Temporal.TaskQueue, worker.Options{})

	app := app.New(
		app.WithConfig(c),
		app.WithWeb(),
		app.WithTemporal(temporal, w),
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
