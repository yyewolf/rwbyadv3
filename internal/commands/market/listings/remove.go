package listings

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/listing"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (cmd *listingsCommand) RemoveListing(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	want := e.SlashCommandInteractionData().Int("card")
	if want < 1 {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Please select a listing number greater than 0...").
			SetEphemeral(true).
			Build(),
		)
	}

	listing, err := p.QueryListings().
		Order(listing.ByCreateTime()).
		Offset(want - 1).
		WithCard().
		First(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("Sorry, you do not have a listing with this number...").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.CommandError(e, err)
	}

	err = utils.WithTx(e.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		listing.Edges.Card.Metadata.Location = "inventory"

		err = tx.Card.UpdateOne(listing.Edges.Card).
			SetAvailable(true).
			SetMetadata(listing.Edges.Card.Metadata).
			Exec(e.Ctx)
		if err != nil {
			return err
		}

		err = tx.Listing.DeleteOne(listing).Exec(e.Ctx)
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
			SetContentf("You have successfully removed the listing for **%s**.", listing.Edges.Card.FullString()).
			SetEphemeral(true),
	)
}
