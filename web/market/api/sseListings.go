package api

import (
	"context"
	"encoding/json"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/models"
)

func (h *MarketApiHandler) OnAddListing(params map[string]interface{}) error {
	var listing = new(models.Listing)
	b, _ := json.Marshal(params["listing"])
	json.Unmarshal(b, listing)

	listing, err := models.Listings(
		qm.Load(
			models.ListingRels.Player,
		),
		qm.Load(
			qm.Rels(models.ListingRels.Card, models.CardRels.CardsStat),
		),
		qm.Where(models.ListingColumns.ID+"=?", listing.ID),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	h.latestListings = append([]*models.Listing{listing}, h.latestListings...)

	return h.SendLatestListings()
}

func (h *MarketApiHandler) OnRemoveListing(params map[string]interface{}) error {
	var listing models.Listing
	b, _ := json.Marshal(params["listing"])
	json.Unmarshal(b, &listing)

	var found bool
	for _, l := range h.latestListings {
		if l.ID == listing.ID {
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
