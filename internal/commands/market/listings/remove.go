package listings

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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
		return utils.CommandError(e, err)
	}

	listing, err := models.Listings(
		qm.Where(models.ListingColumns.CardID+"=?", card.ID),
		qm.Where(models.ListingColumns.PlayerID+"=?", p.ID),
	).One(context.Background(), tx)
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}

	card.Available = true
	utils.Cards.SetLocation(card, "inventory")

	_, err = card.Update(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
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
