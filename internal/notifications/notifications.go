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

	// Possible notifications
	app.Worker().RegisterWorkflow(n.SendDmWorkflow)              // Some component needs to send a DM
	app.Worker().RegisterWorkflow(n.NotifyCardLevelUpWorkflow)   // A card just leveled up
	app.Worker().RegisterWorkflow(n.NotifyPlayerLevelUpWorkflow) // A player just leveled up

	return n
}
