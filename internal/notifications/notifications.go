package notifications

import (
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
)

type NotificationsRepository struct {
	app interfaces.App
}

var Repository *NotificationsRepository

func NewNotificationsRepository(app interfaces.App) *NotificationsRepository {
	eventHandler := app.EventHandler()

	n := &NotificationsRepository{
		app: app,
	}

	Repository = n

	// Events
	eventHandler.OnEvent(jobs.NotifySendDm, n.SendDMJob)

	app.Worker().RegisterWorkflow(n.NotifyCardLevelUpWorkflow)
	// app.Worker().RegisterActivity(cmd.AuctionEndActivity)

	return n
}
