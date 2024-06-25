package auctions

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/temporal"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func (cmd *auctionsCommand) AuctionEndWorkflow(ctx workflow.Context, params *temporal.AuctionEndParams) (*temporal.AuctionEndStatus, error) {
	future := workflow.NewTimer(ctx, params.EndsAt.Sub(workflow.Now(ctx)))

	err := future.Get(ctx, nil)
	if err != nil {
		return &temporal.AuctionEndStatus{
			Status: "error waiting for timer",
		}, err
	}

	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.ID+"=?", params.AuctionID),
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
		return &temporal.AuctionEndStatus{
			Status: "error querying auction",
		}, err
	}

	if auction.EndsAt.After(time.Now()) {
		// reset current workflow
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("end_auction_%s_%d", auction.ID, auction.TimeExtensions),
			TaskQueue: cmd.app.Config().Temporal.TaskQueue,
		}

		params = &temporal.AuctionEndParams{
			AuctionID: params.AuctionID,
			EndsAt:    auction.EndsAt,
		}

		_, err = cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.AuctionEndWorkflow)
		if err != nil {
			return &temporal.AuctionEndStatus{
				Status: "error restarting workflow",
			}, err
		}

		return &temporal.AuctionEndStatus{
			Ok: true,
		}, nil
	}

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var activityResult temporal.AuctionEndStatus
	err = workflow.ExecuteActivity(ctx, cmd.AuctionEndActivity, params.AuctionID).Get(ctx, &activityResult)
	if err != nil {
		return &temporal.AuctionEndStatus{
			Status: "error executing activity",
		}, err
	}

	return &temporal.AuctionEndStatus{
		Ok: true,
	}, nil
}

func (cmd *auctionsCommand) AuctionEndActivity(ctx context.Context, auctionID string) (*temporal.AuctionEndStatus, error) {
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
		return &temporal.AuctionEndStatus{
			Ok: false,
		}, err
	}

	// Check if the auction has any bids
	latestBid, found := utils.Auctions.GetLatestBid(auction)
	if !found {
		// If no bids, give back to seller
		err = cmd.auctionEndNoBid(auction)
		return &temporal.AuctionEndStatus{
			Ok: false,
		}, err
	}

	// If there are bids, give to the highest bidder and give the money back to the other bidders
	err = cmd.auctionEndBidder(auction, latestBid)
	if err != nil {
		return &temporal.AuctionEndStatus{
			Ok: false,
		}, err
	}

	return &temporal.AuctionEndStatus{
		Ok: true,
	}, nil
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
