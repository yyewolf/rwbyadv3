package api

import (
	"context"
	"math"
	"net/url"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/ent/auctionbid"
	"github.com/yyewolf/rwbyadv3/ent/cardtype"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

var (
	auctionsPerPage = 5
)

func (h *MarketApiHandler) fetchAuctionByID(ctx context.Context, auctionID uuid.UUID) (*ent.Auction, error) {
	return h.app.Db().Auction.Query().
		Where(auction.ID(auctionID)).
		Where(auction.EndsAtGTE(time.Now())).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		WithBids(func(abq *ent.AuctionBidQuery) {
			abq.Order(auctionbid.ByPrice(sql.OrderDesc()))
		}).
		Only(ctx)
}

func (h *MarketApiHandler) fetchAuctions(ctx context.Context, offset, limit int, query string) ([]*ent.Auction, error) {
	return h.app.Db().Auction.Query().
		Offset(offset).
		Limit(limit).
		Order(auction.ByCreateTime(sql.OrderDesc())).
		WithOwnedBy(func(pq *ent.PlayerQuery) {
			pq.Where(player.UsernameContains(query))
		}).
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType(func(ctq *ent.CardTypeQuery) {
				ctq.Where(cardtype.NameContains(query))
				ctq.Where(cardtype.SCategoriesContains(query))
			})
		}).
		WithBids(func(abq *ent.AuctionBidQuery) {
			abq.Order(auctionbid.ByPrice(sql.OrderDesc()))
		}).
		Where(auction.EndsAtGTE(time.Now())).
		All(ctx)
}

func handleError(c echo.Context, err error, userMessage string) error {
	c.Response().Header().Add("HX-Retarget", "#message")
	return templates.RenderView(c, market.Error(userMessage))
}

func (h *MarketApiHandler) GetAuctions(c echo.Context) error {
	amount, err := h.app.Db().Auction.Query().
		Where(auction.EndsAtGTE(time.Now())).
		Count(c.Request().Context())
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(c.Request(), auctionsPerPage, amount)

	query := c.QueryParam("q")
	if c.Request().Header.Get("HX-Request") == "true" && query == "" {
		parsedUrl, err := url.ParseRequestURI(c.Request().Header.Get("HX-Current-URL"))
		if err != nil {
			return err
		}
		query = parsedUrl.Query().Get("q")
	}

	auctions, err := h.fetchAuctions(c.Request().Context(), paginator.Offset(), auctionsPerPage, query)
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.Auctions(auctions, paginator))
}

func (h *MarketApiHandler) GetAuction(c echo.Context) error {
	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return err
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.Auction(auction))
}

func (h *MarketApiHandler) GetAuctionPrice(c echo.Context) error {
	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return err
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.AuctionAmount(auction.GetPrice()))
}

func (h *MarketApiHandler) GetAuctionTimeleft(c echo.Context) error {
	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return err
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.AuctionTimeleft(auction))
}

func (h *MarketApiHandler) GetLatestAuctions(c echo.Context) error {
	auctions, err := h.app.Db().Auction.Query().
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
		All(c.Request().Context())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.LatestAuctions(auctions))
}

func (h *MarketApiHandler) GetAuctionModal(c echo.Context) error {
	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return err
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.AuctionModal(auction))
}

func (h *MarketApiHandler) BidOnAuction(c echo.Context) error {
	session := utils.GetSessionFromContext(c)
	bidder := session.Edges.Player

	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return err
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return err
	}

	formBidAmount := c.FormValue("bid")
	bidAmount, err := strconv.ParseInt(formBidAmount, 10, 64)
	if err != nil {
		return handleError(c, err, "You need to enter a number.")
	}

	if bidder.AvailableBalance() < bidAmount {
		// error too poor
		return handleError(c, err, "You do not have enough liens to bid this amount.")
	}

	auctionPrice := auction.GetPrice()

	if bidAmount < auctionPrice+49 {
		// error too poor
		return handleError(c, err, "You need to bid a bit more.")
	}

	// Check for available slots
	if utils.Players.AvailableSlots(bidder) == 0 {
		return handleError(c, err, "You do not have enough slots in your backpack to purchase this card.")
	}

	var latestBid *ent.AuctionBid
	if len(auction.Edges.Bids) > 0 {
		latestBid = auction.Edges.Bids[0]
	}

	// Check if the player is already in the auction
	if latestBid != nil && latestBid.PlayerID == bidder.ID {
		return handleError(c, err, "You are already in the auction.")
	}

	err = ent.WithTx(c.Request().Context(), h.app.Db(), func(tx *ent.Tx) error {
		// Check if there's already a bid in place
		if latestBid != nil {
			err = tx.Player.UpdateOne(latestBid.Edges.Player).
				AddLiens(latestBid.Price).
				AddLiensInAuction(-1 * latestBid.Price).
				AddBackpackReservedSlots(-1).
				Exec(c.Request().Context())
			if err != nil {
				return err
			}
		}

		err = tx.AuctionBid.Create().
			SetPlayer(bidder).
			SetAuction(auction).
			SetPrice(bidAmount).
			Exec(c.Request().Context())
		if err != nil {
			return err
		}

		err = tx.Player.UpdateOne(bidder).
			AddLiensInAuction(bidAmount).
			AddLiens(-1 * bidAmount).
			AddBackpackReservedSlots(1).
			Exec(c.Request().Context())
		if err != nil {
			return err
		}

		delay := 600.0
		if auction.TimeExtensions < 7 {
			delay = delay / (math.Pow(2, float64(auction.TimeExtensions)))
		} else {
			delay = 5
		}

		// Handle time extensions
		if time.Now().After(auction.EndsAt.Add(-time.Duration(delay) * time.Second)) {
			// Calculate new end time
			err = tx.Auction.UpdateOne(auction).
				AddTimeExtensions(1).
				SetEndsAt(auction.EndsAt).
				Exec(c.Request().Context())
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return handleError(c, err, "An error occurred while bidding.")
	}

	cardDescription := auction.Edges.Card.FullString()

	notifications.DispatchDm(h.app, bidder, discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Auction Bid").
				SetColor(h.app.Config().App.BotColor).
				SetDescriptionf("You have bid **%d** Liens on `%s`.", bidAmount, cardDescription).
				SetEmbedFooter(h.app.Footer()).
				Build(),
		).
		Build(),
	)

	if latestBid != nil && latestBid.PlayerID != bidder.ID {
		notifications.DispatchDm(h.app, latestBid.Edges.Player, discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Auction Outbid").
					SetColor(h.app.Config().App.BotColor).
					SetDescriptionf("You have been outbid on `%s` by **%d** Liens.", cardDescription, bidAmount).
					SetEmbedFooter(h.app.Footer()).
					Build(),
			).
			Build(),
		)
	}

	c.Response().Header().Add("HX-Retarget", "#message")
	return templates.RenderView(c, market.Success("You successfully bid on the listing !"))
}
