package events

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/events"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/rates"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
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
			utils.Cards.GiveXP(p.R.SelectedCard, XP)

			workflowOptions := client.StartWorkflowOptions{
				ID:        fmt.Sprintf("card_level_up_%s", e.MessageID),
				TaskQueue: app.Config().Temporal.TaskQueue,
			}
			app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, notifications.Repository.NotifyCardLevelUpWorkflow, &notifications.CardLevelUpParams{
				Player: p,
				Card:   p.R.SelectedCard,
			})

			p.R.SelectedCard.UpdateG(ctx, boil.Infer())
		})
	}
}
