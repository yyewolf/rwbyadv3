package trades

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

const (
	commandName        = "trades"
	commandDescription = "Trades"

	componentId           = "trades/{trade_id}/{action}"
	componentActionAccept = "accept"
	componentActionDeny   = "deny"
)

type tradesCommand struct {
	app interfaces.App
}

func TradesCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd tradesCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/trades/start", builder.WithContext(
				app,
				cmd.Start,
				builder.WithPlayer(),
			))

			// h.ButtonComponent("/"+componentId, builder.WithContextD(
			// 	app,
			// 	cmd.HandleGetListingsInteraction,
			// 	builder.WithPlayer(),
			// ))

			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "start",
					Description: "Start a trade",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "with",
							Description: "Which card do you want to sell ?",
							Required:    true,
						},
					},
				},
			},
		}),
	)
}
