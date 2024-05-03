package stars

import (
	"context"
	"net/url"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "star"
	commandDescription = "Get a reward for starring the github repository"
)

type starCommand struct {
	app interfaces.App
}

func StarCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd starCommand

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

func (cmd *starCommand) HandleCommand(e *handler.CommandEvent) error {
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
	s := p.R.GetGithubStar()
	if s.HasStarred {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("You already starred the repo!").
				SetEphemeral(true),
		)
	}

	state := &models.AuthGithubState{
		State:     uuid.NewString(),
		PlayerID:  p.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Type:      models.AuthGithubStatesTypeCheckStar,
	}
	err := state.InsertG(context.Background(), boil.Infer())
	if err != nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContent("There has been an error").
				SetEphemeral(true),
		)
	}

	cfg := cmd.app.Config()

	url, err := url.JoinPath(cfg.Github.App.BaseURI, "/")
	if err != nil {
		logrus.WithError(err).Fatal("failed to join url path")
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Please star the [repository](https://github.com/%s/%s).\nYou can then click this link to verify your star: %s?s=%s",
				cfg.Github.Username,
				cfg.Github.Repository,
				url,
				state.State,
			).
			SetEphemeral(true),
	)
}
