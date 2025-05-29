package api

import (
	"context"
	"math"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/ent/auctionbid"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/ent/cardtype"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/api"
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

func (h *MarketApiHandler) countAuctions(ctx context.Context, query string) (int, error) {
	return h.app.Db().Auction.Query().
		Where(auction.EndsAtGTE(time.Now())).
		Where(
			auction.Or(
				auction.HasOwnedByWith(
					player.UsernameContainsFold(query),
				),
				auction.HasCardWith(
					card.HasTypeWith(
						cardtype.Or(
							cardtype.NameContainsFold(query),
							cardtype.SCategoriesContainsFold(query),
						),
					),
				),
			),
		).
		Count(ctx)
}

func (h *MarketApiHandler) fetchAuctions(ctx context.Context, offset, limit int, query string) ([]*ent.Auction, error) {
	return h.app.Db().Auction.Query().
		Offset(offset).
		Limit(limit).
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
		Where(
			auction.Or(
				auction.HasOwnedByWith(
					player.UsernameContainsFold(query),
				),
				auction.HasCardWith(
					card.HasTypeWith(
						cardtype.Or(
							cardtype.NameContainsFold(query),
							cardtype.SCategoriesContainsFold(query),
						),
					),
				),
			),
		).
		All(ctx)
}

func (h *MarketApiHandler) GetAuctions(c echo.Context) error {
	amount, err := h.countAuctions(c.Request().Context(), "")
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while counting the auctions.")
	}

	paginator := pagination.NewPaginator(c.Request(), auctionsPerPage, amount)

	query, err := h.GetQueryParam(c, "query")
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while fetching the auctions.")
	}

	auctions, err := h.fetchAuctions(c.Request().Context(), paginator.Offset(), auctionsPerPage, query)
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while fetching the auctions.")
	}

	return api.SendPaginated(c, ent.ViewAuctionListAs(auctions, ent.Public), paginator)
}

func (h *MarketApiHandler) GetAuction(c echo.Context) error {
	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return HandleErrorJson(c, err, "Invalid auction ID.")
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while fetching the auction.")
	}

	return api.SendOK(c, ent.ViewAuctionAs(auction, ent.Public))
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
		return HandleErrorJson(c, err, "An error occurred while fetching latest auctions.")
	}

	return api.SendOK(c, ent.ViewAuctionListAs(auctions, ent.Public))
}

type BidOnAuctionRequest struct {
	BidAmount int64 `json:"bid_amount" validate:"required,min=50"`
}

func (h *MarketApiHandler) BidOnAuction(c echo.Context) error {
	session := utils.GetSessionFromContext(c)
	bidder := session.Edges.Player

	var req BidOnAuctionRequest
	if err := c.Bind(&req); err != nil {
		return HandleErrorJson(c, err, "Invalid request data.")
	}

	auctionID, err := uuid.Parse(c.Param("auctionId"))
	if err != nil {
		return HandleErrorJson(c, err, "Invalid auction ID.")
	}

	auction, err := h.fetchAuctionByID(c.Request().Context(), auctionID)
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while fetching the auction.")
	}

	bidAmount := req.BidAmount

	if bidder.AvailableBalance() < bidAmount {
		// error too poor
		return HandleErrorJson(c, err, "You do not have enough liens to bid this amount.")
	}

	auctionPrice := auction.GetPrice()

	if bidAmount < auctionPrice+49 {
		// error too poor
		return HandleErrorJson(c, err, "You need to bid a bit more.")
	}

	// Check for available slots
	if utils.Players.AvailableSlots(bidder) == 0 {
		return HandleErrorJson(c, err, "You do not have enough slots in your backpack to purchase this card.")
	}

	var latestBid *ent.AuctionBid
	if len(auction.Edges.Bids) > 0 {
		latestBid = auction.Edges.Bids[0]
	}

	// Check if the player is already in the auction
	// if latestBid != nil && latestBid.PlayerID == bidder.ID {
	// 	return HandleErrorJson(c, err, "You are already in the auction.")
	// }

	err = ent.WithTx(c.Request().Context(), h.app.Db(), func(tx *ent.Tx) error {
		// Check if there's already a bid in place
		if latestBid != nil {
			err = tx.Player.UpdateOneID(latestBid.PlayerID).
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
		return HandleErrorJson(c, err, "An error occurred while bidding.")
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

	return api.SendOK(c, true)
}
