package utils

import (
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/models"
)

type auction struct{}

var Auctions auction

func (auction) GetLatestBid(a *models.Auction) (*models.AuctionsBid, bool) {
	if len(a.R.AuctionsBids) == 0 {
		return nil, false
	}

	return a.R.AuctionsBids[0], true
}

func (auction) GetPrice(a *models.Auction) int64 {
	bid, ok := Auctions.GetLatestBid(a)
	if !ok {
		return 1
	}

	return bid.Price
}

func (auction) CreatePlayerBid(a *models.Auction, p *models.Player) *models.AuctionsBid {
	bid := &models.AuctionsBid{
		ID:        uuid.NewString(),
		PlayerID:  p.ID,
		AuctionID: a.ID,
	}

	return bid
}
