package listings

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *listingsCommand) AddListing(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	want := event.SlashCommandInteractionData().Int("card")
	if want < 1 {
		return event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Please select a card number greater than 0...").
			SetEphemeral(true).
			Build(),
		)
	}

	card, err := currentPlayer.QueryCards().
		Where(card.Available(true)).
		Order(card.ByPosition()).
		Offset(want - 1).
		First(event.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sorry, you do not have a card with this number...").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.CommandError(logger, event, err)
	}

	price := int64(event.SlashCommandInteractionData().Int("price"))

	err = ent.WithTx(event.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		_, err := tx.Listing.Create().
			SetPlayerID(currentPlayer.ID).
			SetCardID(card.ID).
			SetPrice(price).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		if currentPlayer.SelectedCardID == card.ID {
			currentPlayer, err = tx.Player.UpdateOne(currentPlayer).
				ClearSelectedCard().
				Save(event.Ctx)
			if err != nil {
				return err
			}
		}

		card.Metadata.Location = "listings"

		_, err = tx.Card.UpdateOne(card).
			SetMetadata(card.Metadata).
			SetAvailable(false).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.Respond(
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
