package events

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/events"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

func OnMessage(app interfaces.App) func(e *events.MessageCreate) {
	return func(e *events.MessageCreate) {
		if e.Message.Author.ID == app.Client().ApplicationID() {
			return
		}

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
			// Don't do XP
			return
		}

		fmt.Println(p.R.SelectedCard)
	}
}
