package api

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
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

	h.Reload()

	return h.SendLatestListings()
}

func (h *MarketApiHandler) SendLatestListings() error {
	var eventData bytes.Buffer
	market.LatestListings(h.latestListings).Render(context.Background(), &eventData)

	h.listeners.Broadcast(&utils.Event{
		Data:  eventData.Bytes(),
		Event: []byte("latest_listings"),
	})

	return nil
}

func (h *MarketApiHandler) SSE(c echo.Context) error {
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	listenerID := uuid.NewString()

	h.listeners.AddNewListener(listenerID, c)

	<-c.Request().Context().Done()
	h.listeners.RemoveListener(listenerID)
	return nil
}
