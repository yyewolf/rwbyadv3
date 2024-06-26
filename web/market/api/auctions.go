package api

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/disgoorg/disgo/discord"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

var (
	auctionsPerPage = 20
)

func (h *MarketApiHandler) GetAuction(c echo.Context) error {
	auctionID := c.Param("auctionId")
	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Where(models.AuctionColumns.ID+"=?", auctionID),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.Auction(auction))
}

func (h *MarketApiHandler) GetAuctionPrice(c echo.Context) error {
	auctionID := c.Param("auctionId")
	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Where(models.AuctionColumns.ID+"=?", auctionID),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.AuctionAmount(utils.Auctions.GetPrice(auction)))
}

func (h *MarketApiHandler) GetAuctionTimeleft(c echo.Context) error {
	auctionID := c.Param("auctionId")
	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Where(models.AuctionColumns.ID+"=?", auctionID),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.AuctionTimeleft(auction))
}

func (h *MarketApiHandler) GetAuctions(c echo.Context) error {
	amount, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt + " > NOW()"),
	).CountG(context.Background())
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(c.Request(), auctionsPerPage, amount)

	auctions, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Offset(paginator.Offset()),
		qm.Limit(auctionsPerPage),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
	).AllG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.Auctions(auctions, paginator))
}

func (h *MarketApiHandler) GetLatestAuctions(c echo.Context) error {
	auctions, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Limit(10),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
		qm.OrderBy(models.AuctionColumns.CreatedAt+" DESC"),
	).AllG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.LatestAuctions(auctions))
}

func (h *MarketApiHandler) GetAuctionModal(c echo.Context) error {
	auctionID := c.Param("auctionId")

	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
		qm.Where(models.AuctionColumns.ID+"=?", auctionID),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.AuctionModal(auction))
}

func (h *MarketApiHandler) BidOnAuction(c echo.Context) error {
	session := utils.GetSessionFromContext(c)
	bidder := session.R.Player

	auctionID := c.Param("auctionId")
	auction, err := models.Auctions(
		qm.Where(models.AuctionColumns.EndsAt+" > NOW()"),
		qm.Load(
			models.AuctionRels.AuctionsBids,
			qm.OrderBy(models.AuctionsBidColumns.Price+" DESC"),
		),
		qm.Load(
			models.AuctionRels.Player,
		),
		qm.Load(
			qm.Rels(models.AuctionRels.Card, models.CardRels.CardsStat),
		),
		qm.Where(models.AuctionColumns.ID+"=?", auctionID),
	).OneG(context.Background())
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}

	formBidAmount := c.FormValue("bid")
	bidAmount, err := strconv.ParseInt(formBidAmount, 10, 64)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You need to enter a number."))
	}

	if utils.Players.AvailableBalance(bidder) < bidAmount {
		// error too poor
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You do not have enough liens to bid this amount."))
	}

	auctionPrice := utils.Auctions.GetPrice(auction)

	if bidAmount < auctionPrice+49 {
		// error too poor
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You need to bid a bit more..."))
	}

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}

	bid := utils.Auctions.CreatePlayerBid(auction, bidder)

	bid.Price = bidAmount
	bidder.LiensBidded += bidAmount
	bidder.Liens -= bidAmount

	// Give back money to latest bidder
	latestBid, found := utils.Auctions.GetLatestBid(auction)
	if found {
		previousBidder, err := models.FindPlayer(context.Background(), tx, latestBid.PlayerID)
		if err != nil {
			tx.Rollback()
			c.Response().Header().Add("HX-Retarget", "#message")
			return templates.RenderView(c, market.Error("An error occured."))
		}

		previousBidder.Liens += latestBid.Price
		previousBidder.LiensBidded -= latestBid.Price
		previousBidder.SlotsReserved--

		_, err = previousBidder.Update(context.Background(), tx, boil.Whitelist(
			models.PlayerColumns.Liens,
			models.PlayerColumns.LiensBidded,
			models.PlayerColumns.SlotsReserved,
		))
		if err != nil {
			tx.Rollback()
			c.Response().Header().Add("HX-Retarget", "#message")
			return templates.RenderView(c, market.Error("An error occured."))
		}
	}

	// Check for available slots
	if utils.Players.AvailableSlots(bidder) == 0 && (latestBid == nil || latestBid.PlayerID != bidder.ID) {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You do not have enough slots in your backpack to purchase this card."))
	}

	bidder.SlotsReserved++

	err = bid.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}

	_, err = bidder.Update(context.Background(), tx, boil.Whitelist(
		models.PlayerColumns.Liens,
		models.PlayerColumns.LiensBidded,
		models.PlayerColumns.SlotsReserved,
	))
	if err != nil {
		tx.Rollback()
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
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
		auction.TimeExtensions++
		auction.EndsAt = auction.EndsAt.Add(time.Duration(delay) * time.Second)

		auction.Update(context.Background(), tx, boil.Whitelist(models.AuctionColumns.EndsAt, models.AuctionColumns.TimeExtensions))
	}

	err = tx.Commit()
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}

	cardDescription := utils.Cards.FullString(auction.R.Card)

	notifications.DispatchDm(h.app, bidder, discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Auction Bid").
				SetColor(h.app.Config().App.BotColor).
				SetDescriptionf("You have bid **%d** Liens on `%s`.", bidAmount, cardDescription).
				Build(),
		).
		Build(),
	)

	if latestBid != nil && latestBid.PlayerID != bidder.ID {
		latestbidder, err := latestBid.Player().OneG(context.Background())
		if err == nil {
			notifications.DispatchDm(h.app, latestbidder, discord.NewMessageCreateBuilder().
				SetEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Auction Outbid").
						SetColor(h.app.Config().App.BotColor).
						SetDescriptionf("You have been outbid on `%s` by **%d** Liens.", cardDescription, bidAmount).
						Build(),
				).
				Build(),
			)
		}
	}

	c.Response().Header().Add("HX-Retarget", "#message")
	return templates.RenderView(c, market.Success("You successfully bid on the listing !"))
}
