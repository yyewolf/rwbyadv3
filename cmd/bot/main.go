package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/app"
	"github.com/yyewolf/rwbyadv3/internal/env"
)

func main() {
	env.Load()
	c := env.Get()

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", c.Database.User, c.Database.Pass, c.Database.Database, c.Database.Host))
	if err != nil {
		logrus.Fatal(err)
	}

	boil.SetDB(db)

	app := app.New(
		app.WithConfig(c),
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
