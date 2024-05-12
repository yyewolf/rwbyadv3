package begin

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "begin"
	commandDescription = "Begin your adventure!"
)

type beginCommand struct {
	app interfaces.App
}

func BeginCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd beginCommand

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

func (cmd *beginCommand) HandleCommand(e *handler.CommandEvent) error {
	tx, err := boil.BeginTx(e.Ctx, nil)
	if err != nil {
		return err
	}

	p := models.Player{
		ID: e.User().ID.String(),
	}
	p.SetGithubStar(context.Background(), tx, true, &models.GithubStar{
		PlayerID: e.User().ID.String(),
	})
	err = p.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("You already have an account!").
				SetEphemeral(true),
		)
	}

	err = tx.Commit()
	if err != nil {
		logrus.WithError(err).WithField("user_id", e.ID().String()).Error("error creating user in db")
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("An error occured :(").
				SetEphemeral(true),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Account created !").
			SetEphemeral(true),
	)
}
