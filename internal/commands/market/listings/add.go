package listings

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

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
	utils.Cards.SetLocation(card, "listings")

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

	tx.Commit()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("All good !").
			SetEphemeral(true),
	)
}
