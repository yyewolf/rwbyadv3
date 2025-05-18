package stars

import (
	"net/url"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

const (
	commandName        = "star"
	commandDescription = "Get a reward for starring the github repository"
)

type starCommand struct {
	app interfaces.App
}

func StarCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd starCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
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

func (cmd *starCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)
	if currentPlayer.Edges.GithubStar.HasStarred {
		return event.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("You already starred the repo!").
				SetEphemeral(true),
		)
	}

	state, err := cmd.app.Db().AuthState.Create().
		SetPlayerID(currentPlayer.ID).
		SetExpiresAt(time.Now().Add(24 * time.Hour)).
		SetType(enums.GithubCheckStar).
		Save(event.Ctx)
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	cfg := cmd.app.Config()

	url, err := url.JoinPath(cfg.Github.App.BaseURI, "/")
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Please star the [repository](https://github.com/%s/%s).\nYou can then click this link to verify your star: %s?s=%s",
				cfg.Github.Username,
				cfg.Github.Repository,
				url,
				state.ID,
			).
			SetEphemeral(true),
	)
}
