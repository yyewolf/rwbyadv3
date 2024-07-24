package listings

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
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
		return utils.CommandError(e, err)
	}

	err = listing.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}

	card.Available = false
	utils.Cards.SetLocation(card, "listings")

	// Remove selected card if it was selected
	if p.SelectedCardID.String == card.ID {
		p.SelectedCardID = null.NewString("", false)
		_, err = p.Update(context.Background(), tx, boil.Whitelist(models.PlayerColumns.SelectedCardID))
		if err != nil {
			tx.Rollback()
			return utils.CommandError(e, err)
		}
	}

	_, err = card.Update(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.CommandError(e, err)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Listings").
					SetDescription("Your listing has been sent !").
					SetColor(cmd.app.Config().App.BotColor).
					SetEmbedFooter(cmd.app.Footer()).
					Build(),
			).
			SetEphemeral(true),
	)
}
