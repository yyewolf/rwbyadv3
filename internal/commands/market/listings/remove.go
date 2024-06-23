package listings

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func (cmd *listingsCommand) RemoveListing(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	want := e.SlashCommandInteractionData().Int("card")
	card, found := utils.Players.GetMarketListing(p, want-1)
	if !found {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, you do not have a card with this number...").
			SetEphemeral(true).
			Build(),
		)
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

	listing, err := models.Listings(
		qm.Where(models.ListingColumns.CardID+"=?", card.ID),
		qm.Where(models.ListingColumns.PlayerID+"=?", p.ID),
	).One(context.Background(), tx)
	if err != nil {
		logrus.WithError(err).Error("could not begin insert listing")
		tx.Rollback()
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Sorry, an error occured.").
			SetEphemeral(true).
			Build(),
		)
	}

	card.Available = true
	utils.Cards.SetLocation(card, "inventory")

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

	listing.Delete(context.Background(), tx, false)

	tx.Commit()

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("All good !").
			SetEphemeral(true),
	)
}
