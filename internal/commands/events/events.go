package events

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/rates"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
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

			p := ctx.Value(builder.PlayerKey).(*models.Player)

			if p.R.SelectedCard == nil {
				// Don't do XP cause no cards are selected
				return
			}

			XP := utils.Cards.GetXPReward(p.R.SelectedCard, 3, false)
			cardLevelUp := utils.Cards.GiveXP(p.R.SelectedCard, XP)
			if cardLevelUp {
				notifications.DispatchCardLevelUp(app, p, p.R.SelectedCard)
			}

			levelBefore := p.Level
			playerLevelUp := utils.Players.GiveXP(p, 1)
			if playerLevelUp {
				notifications.DispatchPlayerLevelUp(app, p, levelBefore)
			}

			p.R.SelectedCard.UpdateG(ctx, boil.Infer())
			p.UpdateG(ctx, boil.Infer())
		})
	}
}
