package listings

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "listings"
	commandDescription = "Listings"
)

type listingsCommand struct {
	app        interfaces.App
	jobHandler interfaces.JobHandler
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
				cmd.AddListing,
				builder.WithPlayer(),
				builder.WithPlayerCards(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			// subcommands for bugs and issues
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "add",
					Description: "Add a listing to the market",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "card",
							Description: "Which card do you want to select ?",
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
					Name:        "remove",
					Description: "Remove a listing from the market",
				},
			},
		}),
	)
}

func (cmd *listingsCommand) AddListing(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	want := e.SlashCommandInteractionData().Int("card")
	card, found := utils.Players.GetAvailableCard(p, want-1)
	if !found {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, you do not have a card with this number...").
			SetEphemeral(true).
			Build(),
		)
	}

	price := int64(e.SlashCommandInteractionData().Int("price"))

	listing := models.Listing{
		ID:       uuid.NewString(),
		PlayerID: p.ID,
		CardID:   card.ID,
		Price:    price,
	}

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		logrus.WithError(err).Error("could not begin tx")
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, an error occured.").
			SetEphemeral(true).
			Build(),
		)
	}

	err = listing.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		logrus.WithError(err).Error("could not begin insert listing")
		tx.Rollback()
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, an error occured.").
			SetEphemeral(true).
			Build(),
		)
	}

	card.Available = false

	_, err = card.Update(context.Background(), tx, boil.Infer())
	if err != nil {
		logrus.WithError(err).Error("could not update card")
		tx.Rollback()
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, an error occured.").
			SetEphemeral(true).
			Build(),
		)
	}

	tx.Commit()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("All good !").
			SetEphemeral(true),
	)
}
