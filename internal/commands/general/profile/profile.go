package profile

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "profile"
	commandDescription = "Displays your profile"
)

type profileCommand struct {
	app interfaces.App
}

func ProfileCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd profileCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerGithubStars(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *profileCommand) HandleCommand(e *handler.CommandEvent) error {
	authError := e.Ctx.Value(builder.ErrorKey)
	if authError != nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("You do not have an account yet...").
				SetEphemeral(true),
		)
	}

	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Loaded : %+#v", p).
			SetEphemeral(true),
	)
}
