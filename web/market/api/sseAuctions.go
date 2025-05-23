package api

import (
	"context"
	"encoding/json"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/ent/auctionbid"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (h *MarketApiHandler) OnAddAuction(params map[string]interface{}) error {
	time.Sleep(100 * time.Millisecond) // Wait for the auction to be inserted into the database
	id := uuid.MustParse(params["id"].(string))

	auction, err := h.app.Db().Auction.Query().
		Where(auction.ID(id)).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		WithBids(func(abq *ent.AuctionBidQuery) {
			abq.Order(auctionbid.ByPrice(sql.OrderDesc()))
		}).
		Only(context.TODO())
	if err != nil {
		return err
	}

	h.latestAuctions = append([]*ent.Auction{auction}, h.latestAuctions...)

	return h.SendLatestAuctions()
}

func (h *MarketApiHandler) OnRemoveAuction(params map[string]interface{}) error {
	id := uuid.MustParse(params["id"].(string))

	h.listeners.Broadcast(&utils.Event{
		Event: []byte("auction_remove"),
		Data:  []byte(id.String()),
	})

	var found bool
	for _, cachedAuction := range h.latestAuctions {
		if cachedAuction.ID == id {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	h.ReloadAuctions()

	return h.SendLatestAuctions()
}

func (h *MarketApiHandler) OnNewBid(params map[string]interface{}) error {
	time.Sleep(100 * time.Millisecond) // Wait for the auction to be inserted into the database
	id := uuid.MustParse(params["id"].(string))

	bid, err := h.app.Db().AuctionBid.Query().
		Where(auctionbid.ID(id)).
		Only(context.TODO())
	if err != nil {
		return err
	}

	eventData, err := json.Marshal(ent.ViewAuctionBidAs(bid, ent.Public))
	if err != nil {
		return err
	}

	h.listeners.Broadcast(&utils.Event{
		Event: []byte("auction_bid"),
		Data:  eventData,
	})

	var found bool
	for _, cachedAuction := range h.latestAuctions {
		if cachedAuction.ID == bid.AuctionID {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	h.ReloadAuctions()

	return nil
}

func (h *MarketApiHandler) OnUpdateAuction(params map[string]interface{}) error {
	id := uuid.MustParse(params["id"].(string))

	var found bool
	for _, cachedAuction := range h.latestAuctions {
		if cachedAuction.ID == id {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	h.ReloadAuctions()

	return h.SendLatestAuctions()
}
