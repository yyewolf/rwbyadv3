package api

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

func (h *MarketApiHandler) SendLatestListings() error {
	var eventData bytes.Buffer
	market.LatestListings(h.latestListings).Render(context.Background(), &eventData)

	h.listeners.Broadcast(&utils.Event{
		Data:  eventData.Bytes(),
		Event: []byte("latest_listings"),
	})

	return nil
}

func (h *MarketApiHandler) SendLatestAuctions() error {
	var eventData bytes.Buffer
	market.LatestAuctions(h.latestAuctions).Render(context.Background(), &eventData)

	h.listeners.Broadcast(&utils.Event{
		Data:  eventData.Bytes(),
		Event: []byte("latest_auctions"),
	})

	return nil
}

func (h *MarketApiHandler) SSE(c echo.Context) error {
	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	listenerID := uuid.NewString()

	h.listeners.AddNewListener(listenerID, c)

	<-c.Request().Context().Done()
	h.listeners.RemoveListener(listenerID)
	return nil
}
