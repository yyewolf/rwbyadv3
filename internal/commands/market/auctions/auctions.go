package auctions

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

const (
	commandName        = "auctions"
	commandDescription = "Auctions"

	componentId            = "auctions/{player_id}/{page}/{action}"
	componentActionPrev    = "prev"
	componentActionRefresh = "refresh"
	componentActionNext    = "next"
)

type auctionsCommand struct {
	app interfaces.App
}

func AuctionsCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd auctionsCommand

	cmd.app = app
	cmd.app.JobHandler().OnEvent(jobs.JobEndAuction, cmd.AuctionEnd)

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/auctions/add", builder.WithContext(
				app,
				cmd.AddAuction,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			h.Command("/auctions/list", builder.WithContext(
				app,
				cmd.GetAuctions,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			h.ButtonComponent("/"+componentId, builder.WithContextD(
				app,
				cmd.HandleGetAuctionsInteraction,
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
							Name:        "duration",
							Description: "How much time should the auction run for ?",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceInt{
								{
									Name:  "1 minute",
									Value: 0,
								},
								{
									Name:  "12 hours",
									Value: 12,
								},
								{
									Name:  "a day",
									Value: 24,
								},
								{
									Name:  "two days",
									Value: 24 * 2,
								},
								{
									Name:  "three days",
									Value: 24 * 3,
								},
								{
									Name:  "a week",
									Value: 24 * 7,
								},
							},
						},
					},
				},
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "list",
					Description: "List all the listings in the market",
				},
			},
		}),
	)
}
