package profile

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
			h.Command("/"+commandName, cmd.HandleCommand)
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *profileCommand) HandleCommand(e *handler.CommandEvent) error {
	p, err := models.FindPlayerG(context.Background(), e.User().ID.String())
	if err != nil {
		p = &models.Player{
			ID: e.User().ID.String(),
		}
		err := p.InsertG(context.Background(), boil.Infer())
		if err != nil {
			return e.Respond(
				discord.InteractionResponseTypeCreateMessage,
				discord.NewMessageCreateBuilder().
					SetContent("There has been an error").
					SetEphemeral(true),
			)
		}

		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("Created : %+#v", p).
				SetEphemeral(true),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Loaded : %+#v", p).
			SetEphemeral(true),
	)
}
