package profile

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

const (
	commandName        = "profile"
	commandDescription = "Displays your profile"
)

type profileCommand struct {
	app interfaces.App
}

func ProfileCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd profileCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
				builder.WithPlayerGithubStars(),
				builder.WithPlayerSelectedCard(),
				builder.WithPlayerLootBoxes(),
				builder.WithPlayerLimits(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *profileCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(cmd.generator(currentPlayer, event.User())),
	)
}
