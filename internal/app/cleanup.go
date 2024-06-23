package app

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/workflow"
)

func (a *App) CleanupJob(ctx workflow.Context) error {
	logrus.WithField("at", time.Now()).Debug("Starting cleanup job")

	var tables = []string{
		models.TableNames.AuthGithubStates,
		models.TableNames.Cards,
		models.TableNames.GithubStars,
		models.TableNames.Jobs,
		models.TableNames.Players,
	}

	for _, table := range tables {
		mods := []qm.QueryMod{
			qm.From(`"` + table + `"`),
			qm.Where(`deleted_at < NOW() - INTERVAL '30 days'`),
		}
		q := models.NewQuery(mods...)
		queries.SetDelete(q)
		q.Exec(boil.GetDB())

		mods = []qm.QueryMod{
			qm.From(`"` + table + `"`),
			qm.Where(`expires_at < NOW()`),
		}
		q = models.NewQuery(mods...)
		queries.SetDelete(q)
		q.Exec(boil.GetDB())
	}

	return nil
}
