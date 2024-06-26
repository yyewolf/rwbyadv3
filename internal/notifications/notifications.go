package notifications

import (
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type NotificationsRepository struct {
	app interfaces.App
}

var Repository *NotificationsRepository

func NewNotificationsRepository(app interfaces.App) *NotificationsRepository {
	n := &NotificationsRepository{
		app: app,
	}

	Repository = n

	// Events
	app.Worker().RegisterWorkflow(n.SendDmWorkflow)
	app.Worker().RegisterWorkflow(n.NotifyCardLevelUpWorkflow)

	return n
}
