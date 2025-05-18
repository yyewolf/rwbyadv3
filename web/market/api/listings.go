package api

import (
	"net/url"

	"entgo.io/ent/dialect/sql"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/cardtype"
	"github.com/yyewolf/rwbyadv3/ent/listing"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

var (
	listingsPerPage = 5
)

func (h *MarketApiHandler) GetListings(c echo.Context) error {
	amount, err := h.app.Db().Listing.Query().
		Count(c.Request().Context())
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(c.Request(), listingsPerPage, amount)

	query := c.QueryParam("q")
	if c.Request().Header.Get("HX-Request") == "true" && query == "" {
		parsedUrl, err := url.ParseRequestURI(c.Request().Header.Get("HX-Current-URL"))
		if err != nil {
			return err
		}
		query = parsedUrl.Query().Get("q")
	}

	listings, err := h.app.Db().Listing.Query().
		Offset(paginator.Offset()).
		Limit(listingsPerPage).
		Order(listing.ByCreateTime(sql.OrderDesc())).
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
		All(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Error("Failed to get listings")
		return err
	}

	return templates.RenderView(c, market.Listings(listings, paginator))
}

func (h *MarketApiHandler) GetLatestListings(c echo.Context) error {
	listings, err := h.app.Db().Listing.Query().
		Limit(10).
		Order(listing.ByCreateTime(sql.OrderDesc())).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		All(c.Request().Context())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.LatestListings(listings))
}

func (h *MarketApiHandler) GetListingModal(c echo.Context) error {
	listingID, err := uuid.Parse(c.Param("listingID"))
	if err != nil {
		return err
	}

	listing, err := h.app.Db().Listing.Query().
		Where(listing.ID(listingID)).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		Only(c.Request().Context())
	if err != nil {
		return err
	}

	return templates.RenderView(c, market.ListingModal(listing))
}

func (h *MarketApiHandler) PurchaseListing(c echo.Context) error {
	session := utils.GetSessionFromContext(c)
	buyer := session.Edges.Player

	listingID, err := uuid.Parse(c.Param("listingID"))
	if err != nil {
		return err
	}

	listing, err := h.app.Db().Listing.Query().
		Where(listing.ID(listingID)).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		Only(c.Request().Context())
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("An error occured."))
	}
	seller := listing.Edges.OwnedBy

	if buyer.AvailableBalance() < listing.Price {
		// error too poor
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You do not have enough liens to purchase this card."))
	}

	// Check for available slots
	if utils.Players.NewAvailableSlots(buyer) == 0 {
		c.Response().Header().Add("HX-Retarget", "#message")
		return templates.RenderView(c, market.Error("You do not have enough slots in your backpack to purchase this card."))
	}

	utils.WithTx(c.Request().Context(), h.app.Db(), func(tx *ent.Tx) error {
		err := tx.Player.UpdateOne(seller).
			AddLiens(listing.Price).
			Exec(c.Request().Context())
		if err != nil {
			return err
		}

		err = tx.Player.UpdateOne(buyer).
			AddLiens(-1 * listing.Price).
			Exec(c.Request().Context())
		if err != nil {
			return err
		}

		listing.Edges.Card.Metadata.Location = "inventory"

		err = tx.Card.UpdateOne(listing.Edges.Card).
			SetPlayerID(buyer.ID).
			SetAvailable(true).
			SetMetadata(listing.Edges.Card.Metadata).
			Exec(c.Request().Context())
		if err != nil {
			return err
		}

		err = tx.Listing.DeleteOne(listing).
			Exec(c.Request().Context())
		if err != nil {
			return err
		}

		return nil
	})

	cardDescription := listing.Edges.Card.FullString()

	notifications.DispatchDm(h.app, buyer, discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Listing Purchase").
				SetColor(h.app.Config().App.BotColor).
				SetDescriptionf("You have purchased `%s` for **%d** Liens.", cardDescription, listing.Price).
				SetEmbedFooter(h.app.Footer()).
				Build(),
		).
		Build(),
	)
	notifications.DispatchDm(h.app, seller, discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Listing Purchase").
				SetColor(h.app.Config().App.BotColor).
				SetDescriptionf("You have sold `%s` for **%d** Liens.", cardDescription, listing.Price).
				SetEmbedFooter(h.app.Footer()).
				Build(),
		).
		Build(),
	)

	c.Response().Header().Add("HX-Retarget", "#message")
	return templates.RenderView(c, market.Success("You successfully purchased the listing !"))
}
