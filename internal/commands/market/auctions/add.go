package auctions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/temporal"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
)

func (cmd *auctionsCommand) AddAuctionB(e *handler.CommandEvent) error {
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

	return cmd.addConfirmation.AskForConfirmation(
		e,
		utils.Joinln(
			fmt.Sprintf("Card: `%s`", utils.Cards.FullString(card)),
			"Are you sure that you want to put this card into auction ? This is **irreversible**.",
		),
		fmt.Sprintf(addConfirmationFormat, want, duration),
	)
}

func (cmd *auctionsCommand) AddAuction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	want, _ := strconv.Atoi(e.Vars["want"])
	card, found := utils.Players.GetAvailableCard(p, want-1)
	if !found {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, you do not have a card with this number...").
			SetEphemeral(true).
			Build(),
		)
	}

	duration, _ := strconv.ParseInt(e.Vars["duration"], 10, 64)

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
		return utils.ComponentError(e, err)
	}

	err = auction.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.ComponentError(e, err)
	}

	card.Available = false
	utils.Cards.SetLocation(card, "auctions")

	// Remove selected card if it was selected
	if p.SelectedCardID.String == card.ID {
		p.SelectedCardID = null.NewString("", false)
		_, err = p.Update(context.Background(), tx, boil.Whitelist(models.PlayerColumns.SelectedCardID))
		if err != nil {
			tx.Rollback()
			return utils.ComponentError(e, err)
		}
	}

	_, err = card.Update(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.ComponentError(e, err)
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
		tx.Rollback()
		return utils.ComponentError(e, err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return utils.ComponentError(e, err)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Auctions").
					SetDescription("Your auction has been sent !").
					SetColor(cmd.app.Config().App.BotColor).
					SetEmbedFooter(cmd.app.Footer()).
					Build(),
			).
			SetEphemeral(true),
	)
}
