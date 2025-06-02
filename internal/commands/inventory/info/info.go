package info

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

const (
	commandName        = "info"
	commandDescription = "Get information about a card in your inventory"
)

type infoCommand struct {
	app interfaces.App
}

func InfoCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd infoCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerSelectedCard(),
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
				},
			},
		}),
	)
}

func (cmd *infoCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	var err error
	var infoCard *ent.Card
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	want := event.SlashCommandInteractionData().Int("card")

	if want >= 1 {
		infoCard, err = currentPlayer.QueryCards().
			Order(card.ByPosition()).
			Where(card.Available(true)).
			WithStats().
			WithType().
			Offset(want - 1).
			First(event.Ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return event.CreateMessage(discord.NewMessageCreateBuilder().
					SetContent("Sorry, you do not have a card with this number...").
					SetEphemeral(true).
					Build(),
				)
			}
			return utils.CommandError(logger, event, err)
		}
	} else {
		if currentPlayer.Edges.SelectedCard == nil {
			return event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("You need to select a card first, use `/select` command.").
				SetEphemeral(true).
				Build(),
			)
		}

		infoCard = currentPlayer.Edges.SelectedCard
	}

	files, embed, component := infoCard.Message()

	message := discord.NewMessageCreateBuilder().
		SetFiles(files).
		AddEmbeds(embed)

	if component != nil {
		message = message.SetContainerComponents(component)
	}

	return event.CreateMessage(message.Build())
}
