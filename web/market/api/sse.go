package api

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

func (h *MarketApiHandler) SendLatestListings() error {
	var eventData bytes.Buffer
	market.LatestListings(h.latestListings).Render(context.Background(), &eventData)

	if eventData.Len() == 0 {
		eventData.WriteByte(' ')
	}

	h.listeners.Broadcast(&utils.Event{
		Data:  eventData.Bytes(),
		Event: []byte("latest_listings"),
	})

	return nil
}

func (h *MarketApiHandler) SendLatestAuctions() error {
	var eventData bytes.Buffer
	market.LatestAuctions(h.latestAuctions).Render(context.Background(), &eventData)

	if eventData.Len() == 0 {
		eventData.WriteByte(' ')
	}

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

	// ping from time to time to keep the connection alive
	go func() {
		ping := func() {
			evt := &utils.Event{
				Data:  []byte("ping"),
				Event: []byte("ping"),
			}
			if err := evt.MarshalTo(w); err != nil {
				h.listeners.RemoveListener(listenerID)
				return
			}
			w.Flush()
		}

		ping()
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ping()
			case <-c.Request().Context().Done():
				return
			}
		}
	}()

	<-c.Request().Context().Done()
	h.listeners.RemoveListener(listenerID)
	return nil
}
