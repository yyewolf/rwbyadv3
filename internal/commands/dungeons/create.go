package dungeons

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
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
			h.Command("/dungeons/create", builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerLimits(),
				builder.WithDungeons(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "create",
					Description: "Create a dungeon!",
				},
			},
		}),
	)
}

func (cmd *beginCommand) HandleCommand(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	// Check if player has an active dungeon
	if len(p.R.Dungeons) > 0 {
		dungeon := p.R.Dungeons[0]
		dungeonUri, err := url.JoinPath(cmd.app.Config().App.BaseURI, "/dungeons/", dungeon.ID)
		if err != nil {
			return utils.CommandError(e, err)
		}

		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				// SetContentf("You had unfinished business ! Join it [here](%s)!", dungeonUri),
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

	if p.R.PlayerLimit.DungeonsLeft <= 0 {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetEphemeral(true).
				SetContentf("You have no more dungeons left!"),
		)
	}

	p.R.PlayerLimit.DungeonsLeft--
	if p.R.PlayerLimit.DungeonsResetAt.IsZero() {
		p.R.PlayerLimit.DungeonsResetAt.SetValid(time.Now().Add(5 * time.Minute))
	}

	tx, err := boil.BeginTx(e.Ctx, nil)
	if err != nil {
		return utils.CommandError(e, err)
	}

	_, err = p.R.PlayerLimit.Update(e.Ctx, tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}

	dungeon := &models.Dungeon{
		ID:       uuid.NewString(),
		PlayerID: p.ID,

		Seed: rand.Int63(),
	}

	err = dungeon.Insert(e.Ctx, tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}

	dungeonUri, err := url.JoinPath(cmd.app.Config().App.BaseURI, "/dungeons/", dungeon.ID)
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}
	tx.Commit()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Dungeon").
					SetDescriptionf("Dungeon created with seed %d, you can join it [here](%s)!", dungeon.Seed, dungeonUri).
					SetColor(cmd.app.Config().App.BotColor).
					SetEmbedFooter(cmd.app.Footer()).
					Build(),
			),
	)
}
