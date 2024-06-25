package auctions

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/temporal"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
)

func (cmd *auctionsCommand) AddAuction(e *handler.CommandEvent) error {
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

	duration := int64(e.SlashCommandInteractionData().Int("duration"))

	auction := models.Auction{
		ID:       uuid.NewString(),
		PlayerID: p.ID,
		CardID:   card.ID,
		EndsAt:   time.Now().Add(time.Duration(duration) * time.Hour),
	}

	if duration == 0 {
		auction.EndsAt = time.Now().Add(time.Minute)
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

	err = auction.Insert(context.Background(), tx, boil.Infer())
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
	utils.Cards.SetLocation(card, "auctions")

	// Remove selected card if it was selected
	if p.SelectedCardID.String == card.ID {
		p.SelectedCardID = null.NewString("", false)
		_, err = p.Update(context.Background(), tx, boil.Whitelist(models.PlayerColumns.SelectedCardID))
		if err != nil {
			logrus.WithError(err).Error("could not update selected card")
			tx.Rollback()
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sorry, an error occured.").
				SetEphemeral(true).
				Build(),
			)
		}
	}

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

	// Schedule end
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("end_auction_%s", auction.ID),
		TaskQueue: cmd.app.Config().Temporal.TaskQueue,
	}

	_, err = cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.AuctionEndWorkflow, &temporal.AuctionEndParams{
		AuctionID: auction.ID,
		EndsAt:    auction.EndsAt,
	})
	if err != nil {
		logrus.WithError(err).Error("failed to schedule delayed end auction job")
		tx.Rollback()
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, an error occured.").
			SetEphemeral(true).
			Build(),
		)
	}

	err = tx.Commit()
	if err != nil {
		logrus.WithError(err).Error("could not update card")
		tx.Rollback()
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, an error occured.").
			SetEphemeral(true).
			Build(),
		)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("All good !").
			SetEphemeral(true),
	)
}
