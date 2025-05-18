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

func OnMessage(app interfaces.App) func(event *events.MessageCreate) {
	return func(event *events.MessageCreate) {
		if event.Message.Author.ID == app.Client().ApplicationID() {
			return
		}

		rateLimiter := getRate(event.Message.Author.ID.String())

		rateLimiter.Do(func() {
			ctx, err := builder.GetContext(
				context.Background(),
				app,
				event.Message.Author.ID,
				builder.WithPlayer(),
				builder.WithPlayerSelectedCard(),
			)
			if err != nil {
				return
			}

			currentPlayer := ctx.Value(builder.NewPlayerKey).(*ent.Player)

			if currentPlayer.Edges.SelectedCard == nil {
				// Don't do XP cause no cards are selected
				return
			}

			experience := currentPlayer.Edges.SelectedCard.GetXPReward(3, false)
			cardLevelUp := currentPlayer.Edges.SelectedCard.GiveXP(experience)
			if cardLevelUp {
				notifications.DispatchCardLevelUp(app, currentPlayer, currentPlayer.Edges.SelectedCard)
			}

			currentPlayer.Edges.SelectedCard.Update().
				SetExperiencePoints(currentPlayer.Edges.SelectedCard.ExperiencePoints).
				SetExperiencePointsThreshold(currentPlayer.Edges.SelectedCard.ExperiencePointsThreshold).
				SetLevel(currentPlayer.Edges.SelectedCard.Level).
				Save(ctx)

			levelBefore := currentPlayer.Level
			playerLevelUp := currentPlayer.GiveXP(1)
			if playerLevelUp {
				notifications.DispatchPlayerLevelUp(app, currentPlayer, levelBefore)
			}

			currentPlayer.Update().
				SetExperiencePoints(currentPlayer.ExperiencePoints).
				SetExperiencePointsThreshold(currentPlayer.ExperiencePointsThreshold).
				SetLevel(currentPlayer.Level).
				Save(ctx)

			// add debug log
			logrus.WithFields(logrus.Fields{
				"author":    event.Message.Author.ID,
				"xp":        experience,
				"level":     currentPlayer.Level,
				"card":      currentPlayer.Edges.SelectedCard.ID,
				"cardxp":    currentPlayer.Edges.SelectedCard.ExperiencePoints,
				"cardlevel": currentPlayer.Edges.SelectedCard.Level,
			}).Debug("semi-passive xp gain")
		})
	}
}
