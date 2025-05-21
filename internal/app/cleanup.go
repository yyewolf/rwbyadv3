package app

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent/authstate"
	"github.com/yyewolf/rwbyadv3/ent/cookie"
	"github.com/yyewolf/rwbyadv3/ent/playerlimit"
	"go.temporal.io/sdk/workflow"
)

func (a *App) CleanupJob(ctx workflow.Context) error {
	logrus.WithField("at", time.Now()).Debug("Starting cleanup job")

	var context = context.TODO()

	a.entClient.Cookie.Delete().
		Where(cookie.ExpiresAtLT(time.Now())).
		Exec(context)

	a.entClient.AuthState.Delete().
		Where(authstate.ExpiresAtLT(time.Now())).
		Exec(context)

	return nil
}

func (a *App) RestoreLimits(ctx workflow.Context) error {
	logrus.WithField("at", time.Now()).Debug("Restore limits")

	var context = context.TODO()

	a.entClient.PlayerLimit.Update().
		Where(playerlimit.DungeonsResetAtLT(time.Now())).
		SetDungeonsResetAt(time.Time{}).
		SetDungeonsLeft(3).
		Exec(context)

	return nil
}
