package selectc

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "select"
	commandDescription = "Select a card to equip it"
)

type selectCommand struct {
	app interfaces.App
}

func SelectCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd selectCommand

	cmd.app = app

	var min = 0
	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "card",
					Description: "Which card do you want to select ?",
					MinValue:    &min,
					Required:    true,
				},
			},
		}),
	)
}

func (cmd *selectCommand) HandleCommand(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	want := e.SlashCommandInteractionData().Int("card")

	if len(p.R.PlayerCards) < want {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, you have to have this card to select it.").
			SetEphemeral(true).
			Build(),
		)
	}

	wantedCard := p.R.PlayerCards[want-1].R.Card
	err := p.SetSelectedCardG(context.Background(), false, wantedCard)
	if err != nil {
		u := uuid.NewString()
		logrus.WithError(err).WithField("error_id", u).Error("a db error occured")
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContentf("Sorry, an error occured... (%s)", u).
			SetEphemeral(true).
			Build(),
		)
	}

	return e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContentf("All good, you selected : `%s`", utils.Cards.FullString(wantedCard)).
		SetEphemeral(true).
		Build(),
	)
}
