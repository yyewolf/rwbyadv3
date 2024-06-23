package jobs

const (
	JobEndAuction  = "end_auction"
	JobCleanupDb   = "cleanup_db"
	JobDelayedPong = "delayed_pong"

	// Web related jobs / notification
	NotifySendDm = "send_dm"

	// Events
	EventNewListing    = "new_listing"
	EventRemoveListing = "remove_listing"
	EventNewAuction    = "new_auction"
	EventRemoveAuction = "remove_auction"
	EventBidAuction    = "bid_auction"
)
