package begin

import (
	"context"
	"strconv"

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

	componentId         = "begin/{player_id}/{page}/{action}"
	componentActionPrev = "prev"
	componentActionNext = "next"
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

			h.ButtonComponent("/"+componentId, builder.WithContextD(
				app,
				cmd.HandleInteraction,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
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

	embed, components := cmd.generator(&p, 0)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *beginCommand) HandleInteraction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	// Get route parameters
	playerID := e.Vars["player_id"]
	action := e.Vars["action"]
	page, _ := strconv.Atoi(e.Vars["page"])

	e.DeferUpdateMessage()
	if playerID != e.User().ID.String() {
		return nil
	}

	switch action {
	case componentActionNext:
		page++
	case componentActionPrev:
		page--
	default:
	}

	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	embed, components := cmd.generator(p, page)

	_, err := e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
