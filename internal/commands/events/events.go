package events

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/rates"
	"golang.org/x/time/rate"
)

func getRate(userId string) *rate.Sometimes {
	val, found := rates.XpRateCache.Get(userId)
	if !found {
		r := &rate.Sometimes{
			Interval: 1 * time.Second,
		}
		rates.XpRateCache.Set(userId, r, 0)
		return r
	}
	return val.(*rate.Sometimes)
}

func OnMessage(app interfaces.App) func(e *events.MessageCreate) {
	return func(e *events.MessageCreate) {
		if e.Message.Author.ID == app.Client().ApplicationID() {
			return
		}

		r := getRate(e.Message.Author.ID.String())

		r.Do(func() {
			ctx, err := builder.GetContext(
				context.Background(),
				app,
				e.Message.Author.ID,
				builder.WithPlayer(),
				builder.WithPlayerSelectedCard(),
			)
			if err != nil {
				return
			}

			p := ctx.Value(builder.NewPlayerKey).(*ent.Player)

			if p.Edges.SelectedCard == nil {
				// Don't do XP cause no cards are selected
				return
			}

			XP := p.Edges.SelectedCard.GetXPReward(3, false)
			cardLevelUp := p.Edges.SelectedCard.GiveXP(XP)
			if cardLevelUp {
				notifications.DispatchCardLevelUp(app, p, p.Edges.SelectedCard)
			}

			p.Edges.SelectedCard.Update().
				SetExperiencePoints(p.Edges.SelectedCard.ExperiencePoints).
				SetExperiencePointsThreshold(p.Edges.SelectedCard.ExperiencePointsThreshold).
				SetLevel(p.Edges.SelectedCard.Level).
				Save(ctx)

			levelBefore := p.Level
			playerLevelUp := p.GiveXP(1)
			if playerLevelUp {
				notifications.DispatchPlayerLevelUp(app, p, levelBefore)
			}

			p.Update().
				SetExperiencePoints(p.ExperiencePoints).
				SetExperiencePointsThreshold(p.ExperiencePointsThreshold).
				SetLevel(p.Level).
				Save(ctx)

			// add debug log
			logrus.WithFields(logrus.Fields{
				"author":    e.Message.Author.ID,
				"xp":        XP,
				"level":     p.Level,
				"card":      p.Edges.SelectedCard.ID,
				"cardxp":    p.Edges.SelectedCard.ExperiencePoints,
				"cardlevel": p.Edges.SelectedCard.Level,
			}).Debug("semi-passive xp gain")
		})
	}
}
