package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func (h *MarketApiHandler) PrepareLatestListings() *utils.Event {
	if len(h.latestListings) == 0 {
		return &utils.Event{
			Data:  []byte("[]"),
			Event: []byte("latest_listings"),
		}
	}

	data, err := json.Marshal(ent.ViewListingListAs(h.latestListings, ent.Public))
	if err != nil {
		return &utils.Event{
			Data:  []byte("error preparing latest listings"),
			Event: []byte("error"),
		}
	}

	return &utils.Event{
		Data:  data,
		Event: []byte("latest_listings"),
	}
}

func (h *MarketApiHandler) PrepareLatestAuctions() *utils.Event {
	if len(h.latestAuctions) == 0 {
		return &utils.Event{
			Data:  []byte("[]"),
			Event: []byte("latest_auctions"),
		}
	}

	data, err := json.Marshal(ent.ViewAuctionListAs(h.latestAuctions, ent.Public))
	if err != nil {
		return &utils.Event{
			Data:  []byte("error preparing latest auctions"),
			Event: []byte("error"),
		}
	}

	return &utils.Event{
		Data:  data,
		Event: []byte("latest_auctions"),
	}
}

func (h *MarketApiHandler) SendLatestListings() error {
	evt := h.PrepareLatestListings()
	if evt == nil {
		return errors.New("failed to prepare latest listings")
	}

	h.listeners.Broadcast(evt)
	return nil
}

func (h *MarketApiHandler) SendLatestAuctions() error {
	evt := h.PrepareLatestAuctions()
	if evt == nil {
		return errors.New("failed to prepare latest auctions")
	}

	h.listeners.Broadcast(evt)
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

		time.Sleep(100 * time.Millisecond) // Give some time for the client to set up the connection
		evt := h.PrepareLatestListings()
		if evt != nil {
			if err := evt.MarshalTo(w); err != nil {
				h.listeners.RemoveListener(listenerID)
				return
			}
			w.Flush()
		}
		evt = h.PrepareLatestAuctions()
		if evt != nil {
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
