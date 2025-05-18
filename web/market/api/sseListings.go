package api

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/listing"
)

func (h *MarketApiHandler) OnAddListing(params map[string]interface{}) error {
	time.Sleep(100 * time.Millisecond) // Wait for the listing to be inserted into the database
	id := uuid.MustParse(params["id"].(string))

	listing, err := h.app.Db().Listing.Query().
		Where(listing.ID(id)).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		Only(context.TODO())
	if err != nil {
		return err
	}

	h.latestListings = append([]*ent.Listing{listing}, h.latestListings...)

	return h.SendLatestListings()
}

func (h *MarketApiHandler) OnRemoveListing(params map[string]interface{}) error {
	time.Sleep(100 * time.Millisecond) // Wait for the listing to be deleted from the database
	id := uuid.MustParse(params["id"].(string))

	var found bool
	for _, cachedListing := range h.latestListings {
		if cachedListing.ID == id {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	h.ReloadListings()

	return h.SendLatestListings()
}
