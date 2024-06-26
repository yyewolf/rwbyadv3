package utils

import (
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/models"
)

type app struct{}

var App app

func (app) DispatchNewListing(app interfaces.App, listing *models.Listing) {
	app.EventHandler().SendEvent(
		jobs.EventNewListing,
		uuid.NewString(),
		map[string]interface{}{
			"listing": listing,
		},
	)
}

func (app) DispatchRemoveListing(app interfaces.App, listing *models.Listing) {
	app.EventHandler().SendEvent(
		jobs.EventRemoveListing,
		uuid.NewString(),
		map[string]interface{}{
			"listing": listing,
		},
	)
}

func (app) DispatchNewAuction(app interfaces.App, auction *models.Auction) {
	app.EventHandler().SendEvent(
		jobs.EventNewAuction,
		uuid.NewString(),
		map[string]interface{}{
			"auction": auction,
		},
	)
}
func (app) DispatchUpdateAuction(app interfaces.App, auction *models.Auction) {
	app.EventHandler().SendEvent(
		jobs.EventUpdateAuction,
		uuid.NewString(),
		map[string]interface{}{
			"auction": auction,
		},
	)
}

func (app) DispatchRemoveAuction(app interfaces.App, auction *models.Auction) {
	app.EventHandler().SendEvent(
		jobs.EventRemoveAuction,
		uuid.NewString(),
		map[string]interface{}{
			"auction": auction,
		},
	)
}

func (app) DispatchNewBid(app interfaces.App, bid *models.AuctionsBid) {
	app.EventHandler().SendEvent(
		jobs.EventBidAuction,
		uuid.NewString(),
		map[string]interface{}{
			"bid": bid,
		},
	)
}
