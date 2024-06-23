package utils

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/models"
)

type A struct{}

var App A

func (A) SendDM(app interfaces.App, to string, msg discord.MessageCreate) {
	app.EventHandler().SendEvent(
		jobs.NotifySendDm,
		uuid.NewString(),
		map[string]interface{}{
			"user_id": to,
			"message": msg,
		},
	)
}

func (A) DispatchNewListing(app interfaces.App, listing *models.Listing) {
	app.EventHandler().SendEvent(
		jobs.EventNewListing,
		uuid.NewString(),
		map[string]interface{}{
			"listing": listing,
		},
	)
}

func (A) DispatchRemoveListing(app interfaces.App, listing *models.Listing) {
	app.EventHandler().SendEvent(
		jobs.EventRemoveListing,
		uuid.NewString(),
		map[string]interface{}{
			"listing": listing,
		},
	)
}

func (A) DispatchNewAuction(app interfaces.App, auction *models.Auction) {
	app.EventHandler().SendEvent(
		jobs.EventNewAuction,
		uuid.NewString(),
		map[string]interface{}{
			"auction": auction,
		},
	)
}

func (A) DispatchRemoveAuction(app interfaces.App, auction *models.Auction) {
	app.EventHandler().SendEvent(
		jobs.EventRemoveAuction,
		uuid.NewString(),
		map[string]interface{}{
			"auction": auction,
		},
	)
}

func (A) DispatchNewBid(app interfaces.App, bid *models.AuctionsBid) {
	app.EventHandler().SendEvent(
		jobs.EventBidAuction,
		uuid.NewString(),
		map[string]interface{}{
			"bid": bid,
		},
	)
}
