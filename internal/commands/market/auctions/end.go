package auctions

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func (cmd *auctionsCommand) AuctionEnd(params map[string]interface{}) error {
	auctionID := params["auction_id"].(string)

	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.ID+"=?", auctionID),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	// Check if the auction has any bids
	latestBid, found := utils.Auctions.GetLatestBid(auction)
	if !found {
		// If no bids, give back to seller
		return cmd.auctionEndNoBid(auction)
	}

	// If there are bids, give to the highest bidder and give the money back to the other bidders
	return cmd.auctionEndBidder(auction, latestBid)
}

func (cmd *auctionsCommand) auctionEndNoBid(auction *models.Auction) error {
	// Give back to seller
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// Card transfer
	card := auction.R.Card
	card.PlayerID = auction.PlayerID
	card.Available = true
	utils.Cards.SetLocation(card, "inventory")

	_, err = card.Update(context.Background(), tx, boil.Whitelist(
		models.CardColumns.PlayerID,
		models.CardColumns.Available,
		models.CardColumns.Metadata,
	))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove auction
	_, err = auction.Delete(context.Background(), tx, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove all auction bidder's
	_, err = models.AuctionsBids(
		qm.Where(models.AuctionsBidColumns.AuctionID+"=?", auction.ID),
	).DeleteAll(context.Background(), tx, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	c := cmd.app.Client()
	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(auction.PlayerID))
	if err != nil {
		return err
	}
	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().SetContentf("Your auction ended and no one bid on it...").Build())
	return err
}

func (cmd *auctionsCommand) auctionEndBidder(auction *models.Auction, latestBid *models.AuctionsBid) error {
	// Give to bidder
	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	seller := auction.R.Player
	bidder, err := models.FindPlayer(context.Background(), tx, latestBid.PlayerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Money tranfer
	bidder.LiensBidded -= latestBid.Price
	seller.Liens += latestBid.Price

	_, err = bidder.Update(context.Background(), tx, boil.Whitelist(
		models.PlayerColumns.LiensBidded,
	))
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = seller.Update(context.Background(), tx, boil.Whitelist(
		models.PlayerColumns.Liens,
	))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Card transfer
	card := auction.R.Card
	card.PlayerID = latestBid.PlayerID
	card.Available = true
	utils.Cards.SetLocation(card, "inventory")

	_, err = card.Update(context.Background(), tx, boil.Whitelist(
		models.CardColumns.PlayerID,
		models.CardColumns.Available,
		models.CardColumns.Metadata,
	))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove auction
	_, err = auction.Delete(context.Background(), tx, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Remove all auction bidder's
	_, err = models.AuctionsBids(
		qm.Where(models.AuctionsBidColumns.AuctionID+"=?", auction.ID),
	).DeleteAll(context.Background(), tx, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// Notify seller
	c := cmd.app.Client()
	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(auction.PlayerID))
	if err != nil {
		return err
	}
	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().SetContentf("Your auction ended and someone got it for : **%d**Ⱡ", latestBid.Price).Build())
	return err
}
