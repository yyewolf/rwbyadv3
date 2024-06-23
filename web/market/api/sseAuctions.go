package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

func (h *MarketApiHandler) OnAddAuction(params map[string]interface{}) error {
	var auction = new(models.Auction)
	b, _ := json.Marshal(params["auction"])
	json.Unmarshal(b, auction)

	auction, err := models.Auctions(
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
		qm.Where(models.AuctionColumns.ID+"=?", auction.ID),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	h.latestAuctions = append([]*models.Auction{auction}, h.latestAuctions...)

	return h.SendLatestAuctions()
}

func (h *MarketApiHandler) OnRemoveAuction(params map[string]interface{}) error {
	var auction models.Auction
	b, _ := json.Marshal(params["auction"])
	json.Unmarshal(b, &auction)

	var found bool
	for _, l := range h.latestAuctions {
		if l.ID == auction.ID {
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
	var bid models.AuctionsBid
	b, _ := json.Marshal(params["bid"])
	json.Unmarshal(b, &bid)

	var eventData bytes.Buffer
	market.AuctionAmount(bid.Price).Render(context.Background(), &eventData)

	h.listeners.Broadcast(&utils.Event{
		Data:  eventData.Bytes(),
		Event: []byte(fmt.Sprintf("auction_%s_bid", bid.AuctionID)),
	})

	return nil
}
