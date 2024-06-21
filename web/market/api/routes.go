package api

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

type MarketApiHandler struct {
	app interfaces.App

	listeners      utils.Listeners
	latestListings []*models.Listing
}

func RegisterAPIRoutes(app interfaces.App, g *echo.Group) {
	handler := MarketApiHandler{app: app}

	app.JobHandler().OnEvent(jobs.EventNewListing, handler.OnAddListing)
	app.JobHandler().OnEvent(jobs.EventRemoveListing, handler.OnRemoveListing)

	handler.Reload()

	// SSE
	g.GET("/sse", handler.SSE)

	// Main page routes
	g.GET("/latest/listings", handler.GetLatestListings)
	g.GET("/latest/auctions", echo.WrapHandler(templ.Handler(market.Main())))

	// Listings routes
	g.GET("/listings", handler.GetListings)
	g.GET("/listings/:listingId", echo.WrapHandler(templ.Handler(market.Main())))

	g.POST("/listings/:listingId", handler.PurchaseListing, auth.DiscordHandler.RequireAuth(discord.WithRedirect("market")))
	g.GET("/listings/:listingId/modal", handler.GetListingModal, auth.DiscordHandler.RequireAuth(discord.WithRedirect("market")))

	// Auctions routes
	g.GET("/auctions", echo.WrapHandler(templ.Handler(market.Main())))
	g.GET("/auctions/:auctionId", echo.WrapHandler(templ.Handler(market.Main())))
	g.POST("/auctions/:auctionId/bid", echo.WrapHandler(templ.Handler(market.Main())))
}

func (h *MarketApiHandler) Reload() {
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
