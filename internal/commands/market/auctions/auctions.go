package auctions

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
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
	cmd.ReconcileAuctions()
	app.EventHandler().OnEvent(jobs.EventRescheduleAuction, cmd.RescheduleAuction)
	app.Worker().RegisterWorkflow(cmd.AuctionEnd)

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

func (cmd *auctionsCommand) ReconcileAuctions() {
	auctions, _ := models.Auctions(
		models.AuctionWhere.EndsAt.LT(time.Now()),
	).AllG(context.Background())

	for i, auction := range auctions {
		// Trigger workflow
		workflowOptions := client.StartWorkflowOptions{
			ID:         fmt.Sprintf("end_auction_%d_%s", auction.TimeExtensions, auction.ID),
			TaskQueue:  "worker",
			StartDelay: time.Duration(i) * time.Second,
		}

		_, err := cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.AuctionEnd, auction.ID)
		if err != nil {
			logrus.WithError(err).Error("failed to schedule delayed end auction job")
		}
	}
}
