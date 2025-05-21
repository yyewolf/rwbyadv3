package inventory

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
	commandName        = "inventory"
	commandDescription = "Check your inventory!"

	componentId            = "inventory/{player_id}/{page}/{action}"
	componentActionPrev    = "prev"
	componentActionRefresh = "refresh"
	componentActionNext    = "next"
)

type inventoryCommand struct {
	app interfaces.App
}

func InventoryCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd inventoryCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
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

func (cmd *inventoryCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := event.User().Username
	if event.User().GlobalName != nil {
		username = *event.User().GlobalName
	}

	embed, components, err := cmd.generator(username, currentPlayer, 0)
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *inventoryCommand) HandleInteraction(logger *logrus.Entry, data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
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

	username := event.User().Username
	if event.User().GlobalName != nil {
		username = *event.User().GlobalName
	}

	embed, components, err := cmd.generator(username, currentPlayer, page)
	if err != nil {
		return utils.ComponentError(logger, event, err)
	}

	_, err = event.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
