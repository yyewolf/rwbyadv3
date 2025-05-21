package dungeons

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/internal/values"
)

const (
	commandName        = "dungeons"
	commandDescription = "All commands for dungeons!"
)

type beginCommand struct {
	app interfaces.App
}

func Command(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd beginCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/dungeons/enter", builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerLimits(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "enter",
					Description: "Enter the dungeon!",
				},
			},
		}),
	)
}

func (cmd *beginCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	dungeon, err := currentPlayer.QueryDungeons().First(event.Ctx)
	if err != nil && !ent.IsNotFound(err) {
		return utils.CommandError(logger, event, err)
	}

	// Check if player has an active dungeon
	if err == nil {
		dungeonUri, err := url.JoinPath(cmd.app.Config().App.BaseURI, "/dungeons/", dungeon.ID.String())
		if err != nil {
			return utils.CommandError(logger, event, err)
		}

		return event.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Dungeon").
						SetDescriptionf("You had unfinished business ! Join it [here](%s)!", dungeonUri).
						SetColor(cmd.app.Config().App.BotColor).
						SetEmbedFooter(cmd.app.Footer()).
						Build(),
				),
		)
	}

	if currentPlayer.Edges.Limits.DungeonsLeft <= 0 {
		return event.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("You have no more dungeons left!"),
		)
	}

	currentPlayer.Edges.Limits.DungeonsLeft--
	if currentPlayer.Edges.Limits.DungeonsResetAt.IsZero() {
		if cmd.app.Config().Mode == values.Prod {
			currentPlayer.Edges.Limits.DungeonsResetAt = time.Now().Add(24 * time.Hour)
		} else {
			currentPlayer.Edges.Limits.DungeonsResetAt = time.Now().Add(5 * time.Minute)
		}
	}

	err = ent.WithTx(event.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		_, err = tx.PlayerLimit.UpdateOne(currentPlayer.Edges.Limits).
			SetDungeonsLeft(currentPlayer.Edges.Limits.DungeonsLeft).
			SetDungeonsResetAt(currentPlayer.Edges.Limits.DungeonsResetAt).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		dungeon, err = tx.Dungeon.Create().
			SetPlayerID(currentPlayer.ID).
			SetSeed(rand.Int63()).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		return nil
	})

	dungeonUri, err := url.JoinPath(cmd.app.Config().App.BaseURI, "/dungeons/", dungeon.ID.String())
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Dungeon").
					SetDescriptionf("Dungeon has been created, you can join it [here](%s)!", dungeonUri).
					SetColor(cmd.app.Config().App.BotColor).
					SetEmbedFooter(cmd.app.Footer()).
					Build(),
			),
	)
}
