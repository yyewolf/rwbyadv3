package api

import (
	"context"
	"net/http"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/ent/auctionbid"
	"github.com/yyewolf/rwbyadv3/ent/listing"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/jobs"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/api"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
)

type MarketApiHandler struct {
	app interfaces.App

	listeners      utils.Listeners
	latestListings []*ent.Listing
	latestAuctions []*ent.Auction
}

func HandleErrorJson(c echo.Context, err error, userMessage string) error {
	logrus.WithError(err).Error(userMessage)
	return api.NewErrorResponse[any](api.ErrorInternalServerError, userMessage).JSON(c, http.StatusInternalServerError)
}

func RegisterAPIRoutes(app interfaces.App, g *echo.Group) {
	handler := MarketApiHandler{app: app}

	app.EventHandler().OnEvent(jobs.EventNewListing, handler.OnAddListing)
	app.EventHandler().OnEvent(jobs.EventRemoveListing, handler.OnRemoveListing)

	app.EventHandler().OnEvent(jobs.EventNewAuction, handler.OnAddAuction)
	app.EventHandler().OnEvent(jobs.EventUpdateAuction, handler.OnUpdateAuction)
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
	g.POST("/listings/:listingId", handler.PurchaseListing, auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectMarket)))

	// Auctions routes
	g.GET("/auctions", handler.GetAuctions)
	g.GET("/auctions/:auctionId", handler.GetAuction)
	g.POST("/auctions/:auctionId", handler.BidOnAuction, auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectMarket)))
}

func (h *MarketApiHandler) ReloadListings() {
	listings, _ := h.app.Db().Listing.Query().
		Limit(10).
		Order(listing.ByCreateTime(sql.OrderDesc())).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		All(context.TODO())

	h.latestListings = listings
}

func (h *MarketApiHandler) ReloadAuctions() {
	auctions, _ := h.app.Db().Auction.Query().
		Limit(10).
		Order(auction.ByCreateTime(sql.OrderDesc())).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		WithBids(func(abq *ent.AuctionBidQuery) {
			abq.Order(auctionbid.ByPrice(sql.OrderDesc()))
		}).
		Where(auction.EndsAtGTE(time.Now())).
		All(context.Background())

	h.latestAuctions = auctions
}
