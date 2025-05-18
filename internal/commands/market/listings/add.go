package listings

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *listingsCommand) AddListing(e *handler.CommandEvent) error {
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

	price := int64(e.SlashCommandInteractionData().Int("price"))

	err = ent.WithTx(e.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		_, err := tx.Listing.Create().
			SetPlayerID(p.ID).
			SetCardID(card.ID).
			SetPrice(price).
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

		card.Metadata.Location = "listings"

		_, err = tx.Card.UpdateOne(card).
			SetMetadata(card.Metadata).
			SetAvailable(false).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		return nil
	})
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
