package begin

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
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

func BeginCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
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
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *beginCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	_, err := cmd.app.Db().Player.Get(event.Ctx, event.User().ID.String())
	if err != nil && !ent.IsNotFound(err) {
		return utils.CommandError(logger, event, err)
	} else if err == nil {
		return event.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetContentf("You already have an account!").
				SetEphemeral(true),
		)
	}

	var player *ent.Player

	err = ent.WithTx(event.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		tempPlayer := ent.Player{
			Level: 1,
		}

		player, err = tx.Player.Create().
			SetID(event.User().ID.String()).
			SetExperiencePointsThreshold(tempPlayer.GetNextLevelXP()).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		_, err = tx.GithubStar.Create().
			SetPlayerID(event.User().ID.String()).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		_, err = tx.PlayerLimit.Create().
			SetPlayerID(event.User().ID.String()).
			SetDungeonsLeft(3).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	embed, components := cmd.generator(player, 0)

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *beginCommand) HandleInteraction(logger *logrus.Entry, data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
	// Get route parameters
	playerID := event.Vars["player_id"]
	action := event.Vars["action"]
	page, _ := strconv.Atoi(event.Vars["page"])

	event.DeferUpdateMessage()
	if playerID != event.User().ID.String() {
		return nil
	}

	switch action {
	case componentActionNext:
		page++
	case componentActionPrev:
		page--
	default:
	}

	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	embed, components := cmd.generator(currentPlayer, page)

	_, err := event.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
