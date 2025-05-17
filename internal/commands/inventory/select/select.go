package selectc

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
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

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
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
					MinValue:    utils.Optional(0),
					Required:    true,
				},
			},
		}),
	)
}

func (cmd *selectCommand) HandleCommand(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	want := e.SlashCommandInteractionData().Int("card")

	if want < 1 {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Please select a card with a number greater than 0...").
			SetEphemeral(true).
			Build(),
		)
	}

	card, err := p.QueryCards().
		Order(card.ByPosition()).
		Where(card.Available(true)).
		Offset(want - 1).
		First(e.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sorry, you do not have a card with this number...").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.CommandError(e, err)
	}

	err = p.Update().SetSelectedCard(card).Exec(e.Ctx)
	if err != nil {
		return utils.CommandError(e, err)
	}

	return e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContentf("All good, you selected : `%s`", card.FullString()).
		SetEphemeral(true).
		Build(),
	)
}
