package auctions

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/ent/auctionbid"
	"github.com/yyewolf/rwbyadv3/internal/temporal"
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

	newCtx := context.Background()

	auction, err := cmd.app.Db().Auction.Query().
		Where(auction.ID(uuid.MustParse(params.AuctionID))).
		Only(newCtx)
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

		_, err = cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.AuctionEndWorkflow, params)
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
	newCtx := context.Background()
	auction, err := cmd.app.Db().Auction.Query().
		Where(auction.ID(uuid.MustParse(auctionID))).
		WithOwnedBy().
		WithBids(func(abq *ent.AuctionBidQuery) {
			abq.WithPlayer()
			abq.Order(auctionbid.ByPrice(sql.OrderDesc()))
		}).
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
		}).
		Only(newCtx)
	if err != nil {
		return &temporal.AuctionEndStatus{
			Ok: false,
		}, err
	}

	if len(auction.Edges.Bids) == 0 {
		// If no bids, give back to seller
		err = cmd.auctionEndNoBid(auction)
		return &temporal.AuctionEndStatus{
			Ok: false,
		}, err
	}

	// Check if the auction has any bids
	latestBid := auction.Edges.Bids[0]

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

func (cmd *auctionsCommand) auctionEndNoBid(auction *ent.Auction) error {
	err := ent.WithTx(context.Background(), cmd.app.Db(), func(tx *ent.Tx) error {
		// Card transfer
		card := auction.Edges.Card
		card.PlayerID = auction.PlayerID
		card.Available = true
		card.Metadata.Location = "inventory"

		card, err := tx.Card.UpdateOne(card).
			SetPlayerID(card.PlayerID).
			SetAvailable(card.Available).
			SetMetadata(card.Metadata).
			Save(context.Background())
		if err != nil {
			return err
		}

		_, err = tx.AuctionBid.Delete().Where(auctionbid.AuctionID(auction.ID)).Exec(context.Background())
		if err != nil {
			return err
		}

		err = tx.Auction.DeleteOne(auction).Exec(context.Background())
		if err != nil {
			return err
		}

		return nil
	})
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

func (cmd *auctionsCommand) auctionEndBidder(auction *ent.Auction, latestBid *ent.AuctionBid) error {
	err := ent.WithTx(context.Background(), cmd.app.Db(), func(tx *ent.Tx) error {
		seller := auction.Edges.OwnedBy
		bidder := latestBid.Edges.Player

		// Money transfer
		bidder, err := tx.Player.UpdateOne(bidder).
			AddLiensInAuction(-1 * latestBid.Price).
			AddBackpackReservedSlots(-1).
			Save(context.Background())
		if err != nil {
			return err
		}

		seller, err = tx.Player.UpdateOne(seller).
			AddLiensInAuction(latestBid.Price).
			Save(context.Background())
		if err != nil {
			return err
		}

		// Card transfer
		card := auction.Edges.Card
		card.Metadata.Location = "inventory"

		card, err = tx.Card.UpdateOne(card).
			SetPlayerID(latestBid.PlayerID).
			SetAvailable(true).
			SetMetadata(card.Metadata).
			Save(context.Background())
		if err != nil {
			return err
		}

		_, err = tx.AuctionBid.Delete().Where(auctionbid.AuctionID(auction.ID)).Exec(context.Background())
		if err != nil {
			return err
		}

		err = tx.Auction.DeleteOne(auction).Exec(context.Background())
		if err != nil {
			return err
		}

		return nil
	})
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
