package utils

import (
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
)

type app struct{}

var App app

func (app) DispatchNewListing(app interfaces.App, id uuid.UUID) {
	app.EventHandler().SendEvent(
		jobs.EventNewListing,
		uuid.NewString(),
		map[string]interface{}{
			"id": id,
		},
	)
}

func (app) DispatchRemoveListing(app interfaces.App, id uuid.UUID) {
	app.EventHandler().SendEvent(
		jobs.EventRemoveListing,
		uuid.NewString(),
		map[string]interface{}{
			"id": id,
		},
	)
}

func (app) DispatchNewAuction(app interfaces.App, id uuid.UUID) {
	app.EventHandler().SendEvent(
		jobs.EventNewAuction,
		uuid.NewString(),
		map[string]interface{}{
			"id": id,
		},
	)
}
func (app) DispatchUpdateAuction(app interfaces.App, id uuid.UUID) {
	app.EventHandler().SendEvent(
		jobs.EventUpdateAuction,
		uuid.NewString(),
		map[string]interface{}{
			"id": id,
		},
	)
}

func (app) DispatchRemoveAuction(app interfaces.App, id uuid.UUID) {
	app.EventHandler().SendEvent(
		jobs.EventRemoveAuction,
		uuid.NewString(),
		map[string]interface{}{
			"id": id,
		},
	)
}

func (app) DispatchNewBid(app interfaces.App, id uuid.UUID) {
	app.EventHandler().SendEvent(
		jobs.EventBidAuction,
		uuid.NewString(),
		map[string]interface{}{
			"id": id,
		},
	)
}
