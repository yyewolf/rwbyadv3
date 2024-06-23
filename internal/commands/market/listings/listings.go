package listings

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

const (
	commandName        = "listings"
	commandDescription = "Listings"

	componentId            = "listings/{player_id}/{page}/{action}"
	componentActionPrev    = "prev"
	componentActionRefresh = "refresh"
	componentActionNext    = "next"
)

type listingsCommand struct {
	app interfaces.App
}

func ListingsCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd listingsCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/listings/add", builder.WithContext(
				app,
				cmd.AddListing,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			h.Command("/listings/remove", builder.WithContext(
				app,
				cmd.RemoveListing,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))

			h.Command("/listings/list", builder.WithContext(
				app,
				cmd.GetListings,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			h.ButtonComponent("/"+componentId, builder.WithContextD(
				app,
				cmd.HandleGetListingsInteraction,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "add",
					Description: "Add a listing to the market",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "card",
							Description: "Which card do you want to sell ?",
							MinValue:    utils.Optional(0),
							Required:    true,
						},
						discord.ApplicationCommandOptionInt{
							Name:        "price",
							Description: "At what price do you wish to sell it ?",
							MinValue:    utils.Optional(1),
							Required:    true,
						},
					},
				},
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "list",
					Description: "List all the listings in the market",
				},
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "remove",
					Description: "Remove a listing from the market",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "card",
							Description: "Which card do you want to remove from the market ?",
							MinValue:    utils.Optional(0),
							Required:    true,
						},
					},
				},
			},
		}),
	)
}
