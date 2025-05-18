package dungeons

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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

func Command(ms *builder.MenuStore, app interfaces.App) *builder.Command {
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

func (cmd *beginCommand) HandleCommand(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	dungeon, err := p.QueryDungeons().First(e.Ctx)
	if err != nil && !ent.IsNotFound(err) {
		return utils.CommandError(e, err)
	}

	// Check if player has an active dungeon
	if err == nil {
		dungeonUri, err := url.JoinPath(cmd.app.Config().App.BaseURI, "/dungeons/", dungeon.ID.String())
		if err != nil {
			return utils.CommandError(e, err)
		}

		return e.Respond(
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

	if p.Edges.Limits.DungeonsLeft <= 0 {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("You have no more dungeons left!"),
		)
	}

	p.Edges.Limits.DungeonsLeft--
	if p.Edges.Limits.DungeonsResetAt.IsZero() {
		if cmd.app.Config().Mode == values.Prod {
			p.Edges.Limits.DungeonsResetAt = time.Now().Add(24 * time.Hour)
		} else {
			p.Edges.Limits.DungeonsResetAt = time.Now().Add(5 * time.Minute)
		}
	}

	err = ent.WithTx(e.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		_, err = tx.PlayerLimit.UpdateOne(p.Edges.Limits).
			SetDungeonsLeft(p.Edges.Limits.DungeonsLeft).
			SetDungeonsResetAt(p.Edges.Limits.DungeonsResetAt).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		dungeon, err = tx.Dungeon.Create().
			SetPlayerID(p.ID).
			SetSeed(rand.Int63()).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		return nil
	})

	dungeonUri, err := url.JoinPath(cmd.app.Config().App.BaseURI, "/dungeons/", dungeon.ID.String())
	if err != nil {
		return utils.CommandError(e, err)
	}

	return e.Respond(
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
