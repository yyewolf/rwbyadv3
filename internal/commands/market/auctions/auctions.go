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
	"github.com/yyewolf/rwbyadv3/internal/temporal"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/internal/utils/confirmation"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
)

const (
	commandName        = "auctions"
	commandDescription = "Auctions"

	addConfirmationId     = "auctions/add/{want}/{duration}/confirmation"
	addConfirmationFormat = "auctions/add/%d/%d/confirmation"

	componentId            = "auctions/list/{player_id}/{page}/{action}"
	componentFormat        = "auctions/list/%s/%d/%s"
	componentActionPrev    = "prev"
	componentActionRefresh = "refresh"
	componentActionNext    = "next"
)

type auctionsCommand struct {
	app interfaces.App

	addConfirmation *confirmation.Handler
}

func AuctionsCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd auctionsCommand

	cmd.app = app
	cmd.ReconcileAuctions()
	app.Worker().RegisterWorkflow(cmd.AuctionEndWorkflow)
	app.Worker().RegisterActivity(cmd.AuctionEndActivity)

	cmd.addConfirmation = confirmation.NewHandler(app, addConfirmationId, builder.WithContextD(
		app,
		cmd.AddAuction,
		builder.WithPlayer(),
		builder.WithPlayerCards(),
	))

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/auctions/add", builder.WithContext(
				app,
				cmd.AddAuctionB,
				builder.WithPlayer(),
			))

			h.Command("/auctions/list", builder.WithContext(
				app,
				cmd.GetAuctions,
				builder.WithPlayer(),
			))

			h.ButtonComponent("/"+componentId, builder.WithContextD(
				app,
				cmd.HandleGetAuctionsInteraction,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))

			cmd.addConfirmation.SetupMux(h)

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
	auctions, _ := models.Auctions().AllG(context.Background())

	for i, auction := range auctions {
		// Trigger workflow
		workflowOptions := client.StartWorkflowOptions{
			ID:         fmt.Sprintf("end_auction_%s", auction.ID),
			TaskQueue:  cmd.app.Config().Temporal.TaskQueue,
			StartDelay: time.Duration(i) * time.Second,
		}

		if auction.TimeExtensions > 0 {
			workflowOptions.ID = fmt.Sprintf("end_auction_%s_%d", auction.ID, auction.TimeExtensions)
		}

		_, err := cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.AuctionEndWorkflow, &temporal.AuctionEndParams{
			AuctionID: auction.ID,
			EndsAt:    auction.EndsAt,
		})
		if err != nil {
			logrus.WithError(err).Error("failed to schedule delayed end auction job")
		}
	}
}
