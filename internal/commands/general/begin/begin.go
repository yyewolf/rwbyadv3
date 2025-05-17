package begin

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
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
			h.Command("/"+commandName, builder.WithContext(
				cmd.app,
				cmd.HandleCommand,
			))

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
	_, err := cmd.app.Db().Player.Get(e.Ctx, e.User().ID.String())
	if err != nil && !ent.IsNotFound(err) {
		return utils.CommandError(e, err)
	} else if err == nil {
		return e.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("You already have an account!").
				SetEphemeral(true),
		)
	}

	var player *ent.Player

	err = utils.WithTx(e.Ctx, cmd.app.Db(), func(ntx *ent.Tx) error {
		tempPlayer := ent.Player{
			Level: 1,
		}

		player, err = ntx.Player.Create().
			SetID(e.User().ID.String()).
			SetExperiencePointsThreshold(tempPlayer.GetNextLevelXP()).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		_, err = ntx.GithubStar.Create().
			SetPlayerID(e.User().ID.String()).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		_, err = ntx.PlayerLimit.Create().
			SetPlayerID(e.User().ID.String()).
			SetDungeonsLeft(3).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return utils.CommandError(e, err)
	}

	embed, components := cmd.generator(player, 0)

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

	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	embed, components := cmd.generator(p, page)

	_, err := e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
