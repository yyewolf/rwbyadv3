package app

import (
	"context"
	"encoding/json"
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

	b, _ := json.Marshal(models.TableNames)
	var tables = make(map[string]string)
	json.Unmarshal(b, &tables)

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

func (a *App) RestoreLimits(ctx workflow.Context) error {
	logrus.WithField("at", time.Now()).Debug("Restore limits")

	// Restore player limits
	models.PlayerLimits(
		qm.Where(models.PlayerLimitColumns.DungeonsResetAt+" < NOW()"),
	).UpdateAllG(
		context.Background(),
		models.M{
			models.PlayerLimitColumns.DungeonsResetAt: nil,
			models.PlayerLimitColumns.DungeonsLeft:    3,
		},
	)

	return nil
}
