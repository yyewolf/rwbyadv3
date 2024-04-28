package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/app"
	"github.com/yyewolf/rwbyadv3/internal/database"
	"github.com/yyewolf/rwbyadv3/internal/env"
)

func main() {
	env.Load()
	c := env.Get()

	db := database.New(c)
	db.Migrate()

	app := app.New(
		app.WithConfig(c),
		app.WithDatabase(db),
	)

	go app.Start()

	// Listen for CTRL+C
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	logrus.Info("Bot is now running. Press CTRL+C to exit.")
	<-done // Will block here until user hits ctrl+c

	app.Shutdown()
}
