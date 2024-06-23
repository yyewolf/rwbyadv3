package api

import (
	"context"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/disgoorg/disgo/discord"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

var (
	listingsPerPage = 20
)

func (h *MarketApiHandler) GetListings(c echo.Context) error {
	amount, err := models.Listings().CountG(context.Background())
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(c.Request(), listingsPerPage, amount)

	listings, err := models.Listings(
		qm.Offset(paginator.Offset()),
		qm.Limit(listingsPerPage),
		qm.Load(
			models.ListingRels.Player,
		),
		qm.Load(
			qm.Rels(models.ListingRels.Card, models.CardRels.CardsStat),
		),
	).AllG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.Listings(listings, paginator))
}

func (h *MarketApiHandler) GetLatestListings(c echo.Context) error {
	listings, err := models.Listings(
		qm.Limit(10),
		qm.Load(
			models.ListingRels.Player,
		),
		qm.Load(
			qm.Rels(models.ListingRels.Card, models.CardRels.CardsStat),
		),
		qm.OrderBy(models.ListingColumns.CreatedAt+" DESC"),
	).AllG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.LatestListings(listings))
}

func (h *MarketApiHandler) GetListingModal(c echo.Context) error {
	listingID := c.Param("listingId")

	listing, err := models.Listings(
		qm.Load(
			models.ListingRels.Player,
		),
		qm.Load(
			qm.Rels(models.ListingRels.Card, models.CardRels.CardsStat),
		),
		qm.Where(models.ListingColumns.ID+"=?", listingID),
	).OneG(context.Background())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.ListingModal(listing))
}

func (h *MarketApiHandler) PurchaseListing(c echo.Context) error {
	session := utils.GetSessionFromContext(c)
	buyer := session.R.Player

	listingID := c.Param("listingId")
	listing, err := models.Listings(
		qm.Load(
			models.ListingRels.Player,
		),
		qm.Load(
			qm.Rels(models.ListingRels.Card, models.CardRels.CardsStat),
		),
		qm.Where(models.ListingColumns.ID+"=?", listingID),
	).OneG(context.Background())
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}
	seller := listing.R.Player

	if utils.Players.AvailableBalance(buyer) < listing.Price {
		// error too poor
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You do not have enough liens to purchase this card."))
	}

	// TODO: Add max card check

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}

	// Money tranfer
	if seller.ID != buyer.ID {
		seller.Liens += listing.Price
		buyer.Liens -= listing.Price
	}

	seller.Update(context.Background(), tx, boil.Whitelist(
		models.PlayerColumns.Liens,
	))
	buyer.Update(context.Background(), tx, boil.Whitelist(
		models.PlayerColumns.Liens,
	))

	// Card transfer
	card := listing.R.Card
	card.PlayerID = buyer.ID
	card.Available = true
	utils.Cards.SetLocation(card, "inventory")

	card.Update(context.Background(), tx, boil.Whitelist(
		models.CardColumns.PlayerID,
		models.CardColumns.Available,
		models.CardColumns.Metadata,
	))

	// Remove listing
	listing.Delete(context.Background(), tx, false)

	tx.Commit()

	cardDescription := utils.Cards.FullString(listing.R.Card)

	utils.App.SendDM(h.app, buyer.ID, discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Listing Purchase").
				SetColor(h.app.Config().App.BotColor).
				SetDescriptionf("You have purchased `%s` for **%d** Liens.", cardDescription, listing.Price).
				Build(),
		).
		Build(),
	)
	utils.App.SendDM(h.app, seller.ID, discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Listing Purchase").
				SetColor(h.app.Config().App.BotColor).
				SetDescriptionf("You have sold `%s` for **%d** Liens.", cardDescription, listing.Price).
				Build(),
		).
		Build(),
	)

	c.Response().Header().Add("HX-Retarget", "#message")
	return templates.RenderView(c, market.Success("You successfully purchased the listing !"))
}
