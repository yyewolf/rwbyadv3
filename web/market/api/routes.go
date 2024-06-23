package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
)

type MarketApiHandler struct {
	app interfaces.App

	listeners      utils.Listeners
	latestListings []*models.Listing
	latestAuctions []*models.Auction
}

func RegisterAPIRoutes(app interfaces.App, g *echo.Group) {
	handler := MarketApiHandler{app: app}

	app.EventHandler().OnEvent(jobs.EventNewListing, handler.OnAddListing)
	app.EventHandler().OnEvent(jobs.EventRemoveListing, handler.OnRemoveListing)
	app.EventHandler().OnEvent(jobs.EventNewAuction, handler.OnAddAuction)
	app.EventHandler().OnEvent(jobs.EventRemoveAuction, handler.OnRemoveAuction)
	app.EventHandler().OnEvent(jobs.EventBidAuction, handler.OnNewBid)

	handler.ReloadListings()
	handler.ReloadAuctions()

	// SSE
	g.GET("/sse", handler.SSE)

	// Main page routes
	g.GET("/latest/listings", handler.GetLatestListings)
	g.GET("/latest/auctions", handler.GetLatestAuctions)

	// Listings routes
	g.GET("/listings", handler.GetListings)
	// g.GET("/listings/:listingId", echo.WrapHandler(templ.Handler(market.Main()))) Not required, maybe later :D
	g.POST("/listings/:listingId", handler.PurchaseListing, auth.DiscordHandler.RequireAuth(discord.WithRedirect("market")))
	g.GET("/listings/:listingId/modal", handler.GetListingModal, auth.DiscordHandler.RequireAuth(discord.WithRedirect("market")))

	// Auctions routes
	g.GET("/auctions", handler.GetAuctions)
	g.GET("/auctions/:auctionId", handler.GetAuction)
	g.GET("/auctions/:auctionId/price", handler.GetAuctionPrice)
	g.POST("/auctions/:auctionId", handler.BidOnAuction, auth.DiscordHandler.RequireAuth(discord.WithRedirect("market")))
	g.GET("/auctions/:auctionId/modal", handler.GetAuctionModal, auth.DiscordHandler.RequireAuth(discord.WithRedirect("market")))
}

func (h *MarketApiHandler) ReloadListings() {
	listings, _ := models.Listings(
		qm.Limit(10),
		qm.Load(
			models.ListingRels.Player,
		),
		qm.Load(
			qm.Rels(models.ListingRels.Card, models.CardRels.CardsStat),
		),
		qm.OrderBy(models.ListingColumns.CreatedAt+" DESC"),
	).AllG(context.Background())

	h.latestListings = listings
}

func (h *MarketApiHandler) ReloadAuctions() {
	auctions, _ := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Limit(10),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
		qm.OrderBy(models.AuctionColumns.CreatedAt+" DESC"),
	).AllG(context.Background())

	h.latestAuctions = auctions
}
