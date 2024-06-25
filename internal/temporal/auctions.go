package temporal

import "time"

type AuctionEndParams struct {
	AuctionID string
	EndsAt    time.Time
}

type AuctionEndStatus struct {
	Ok     bool
	Status string
}
