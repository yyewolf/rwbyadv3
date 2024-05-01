package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/app"
	"github.com/yyewolf/rwbyadv3/internal/env"
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

	app := app.New(
		app.WithConfig(c),
		app.WithWeb(),
	)

	go app.Start()

	// Listen for CTRL+C
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	logrus.Info("Bot is now running. Press CTRL+C to exit.")
	<-done // Will block here until user hits ctrl+c

	db.Close()
	app.Shutdown()
}
