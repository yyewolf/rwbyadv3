package auctions

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/temporal"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"go.temporal.io/sdk/client"
)

func (cmd *auctionsCommand) AddAuctionB(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	want := e.SlashCommandInteractionData().Int("card")
	if want < 1 {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Please select a card number greater than 0...").
			SetEphemeral(true).
			Build(),
		)
	}

	card, err := p.QueryCards().
		Where(card.Available(true)).
		Order(card.ByPosition()).
		Offset(want - 1).
		First(e.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sorry, you do not have a card with this number...").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.CommandError(e, err)
	}

	duration := int64(e.SlashCommandInteractionData().Int("duration"))

	return cmd.addConfirmation.AskForConfirmation(
		e,
		utils.Joinln(
			fmt.Sprintf("Card: `%s`", card.FullString()),
			"Are you sure that you want to put this card into auction ? This is **irreversible**.",
		),
		fmt.Sprintf(addConfirmationFormat, want, duration),
	)
}

func (cmd *auctionsCommand) AddAuction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	want, _ := strconv.Atoi(e.Vars["want"])
	if want < 1 {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Please select a card number greater than 0...").
			SetEphemeral(true).
			Build(),
		)
	}

	card, err := p.QueryCards().
		Where(card.Available(true)).
		Order(card.ByPosition()).
		Offset(want - 1).
		First(e.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sorry, you do not have a card with this number...").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.ComponentError(e, err)
	}

	duration, _ := strconv.ParseInt(e.Vars["duration"], 10, 64)

	err = utils.WithTx(e.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		endsAt := time.Now().Add(time.Duration(duration) * time.Hour)
		if duration == 0 {
			endsAt = time.Now().Add(time.Minute)
		}

		auction, err := tx.Auction.Create().
			SetPlayerID(p.ID).
			SetCardID(card.ID).
			SetEndsAt(endsAt).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		if p.SelectedCardID == card.ID {
			p, err = tx.Player.UpdateOne(p).
				ClearSelectedCard().
				Save(e.Ctx)
			if err != nil {
				return err
			}
		}

		card.Metadata.Location = "auctions"

		_, err = tx.Card.UpdateOne(card).
			SetMetadata(card.Metadata).
			SetAvailable(false).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		// Schedule end
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("end_auction_%s", auction.ID),
			TaskQueue: cmd.app.Config().Temporal.TaskQueue,
		}

		_, err = cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.AuctionEndWorkflow, &temporal.AuctionEndParams{
			AuctionID: auction.ID.String(),
			EndsAt:    auction.EndsAt,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
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
