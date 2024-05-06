package inventory

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
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

func InventoryCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd inventoryCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))

			h.ButtonComponent("/"+componentId, builder.WithContextD(
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

func (cmd *inventoryCommand) HandleCommand(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	username := e.User().Username
	if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	}

	embed, components := cmd.generator(username, p, 0)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *inventoryCommand) HandleInteraction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
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

	username := e.User().Username
	if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	}

	embed, components := cmd.generator(username, p, page)

	_, err := e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
