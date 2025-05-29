package api

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/ent/cardtype"
	"github.com/yyewolf/rwbyadv3/ent/listing"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/api"
)

var (
	listingsPerPage = 5
)

func (h *MarketApiHandler) fetchListingByID(ctx context.Context, listingID uuid.UUID) (*ent.Listing, error) {
	return h.app.Db().Listing.Query().
		Where(listing.ID(listingID)).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		Only(ctx)
}

func (h *MarketApiHandler) countListings(ctx context.Context, query string) (int, error) {
	return h.app.Db().Listing.Query().
		Where(
			listing.Or(
				listing.HasOwnedByWith(
					player.UsernameContainsFold(query),
				),
				listing.HasCardWith(
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

func (h *MarketApiHandler) fetchListings(ctx context.Context, offset, limit int, query string) ([]*ent.Listing, error) {
	return h.app.Db().Listing.Query().
		Offset(offset).
		Limit(listingsPerPage).
		Order(listing.ByCreateTime(sql.OrderDesc())).
		WithOwnedBy().
		WithCard(func(cq *ent.CardQuery) {
			cq.WithStats()
			cq.WithType()
		}).
		Where(
			listing.Or(
				listing.HasOwnedByWith(
					player.UsernameContainsFold(query),
				),
				listing.HasCardWith(
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

func (h *MarketApiHandler) GetQueryParam(c echo.Context, param string) (string, error) {
	query := c.QueryParam(param)
	return query, nil
}

func (h *MarketApiHandler) GetListings(c echo.Context) error {
	amount, err := h.countListings(c.Request().Context(), c.QueryParam("query"))
	if err != nil {
		return HandleErrorJson(c, fmt.Errorf("could not list amount: %w", err), "An error occurred while fetching listings.")
	}

	paginator := pagination.NewPaginator(c.Request(), listingsPerPage, amount)

	query := c.QueryParam("query")
	if err != nil {
		return HandleErrorJson(c, fmt.Errorf("could not get query param: %w", err), "An error occurred while fetching listings.")
	}

	listings, err := h.fetchListings(c.Request().Context(), 0, listingsPerPage, query)
	if err != nil {
		return HandleErrorJson(c, fmt.Errorf("could not fetch listings: %w", err), "An error occurred while fetching listings.")
	}

	return api.SendPaginated(c, ent.ViewListingListAs(listings, ent.Public), paginator)
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
		return HandleErrorJson(c, err, "An error occurred while fetching latest listings.")
	}

	return api.SendOK(c, ent.ViewListingListAs(listings, ent.Public))
}

func (h *MarketApiHandler) PurchaseListing(c echo.Context) error {
	session := utils.GetSessionFromContext(c)
	buyer := session.Edges.Player

	listingID, err := uuid.Parse(c.Param("listingId"))
	if err != nil {
		return HandleErrorJson(c, err, "Invalid listing ID.")
	}

	listing, err := h.fetchListingByID(c.Request().Context(), listingID)
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while fetching the listing.")
	}
	seller := listing.Edges.OwnedBy

	if buyer.AvailableBalance() < listing.Price {
		return api.SendValidationError(c, "You do not have enough liens to purchase this card.", nil)
	}

	// Check for available slots
	if utils.Players.AvailableSlots(buyer) == 0 {
		return api.SendValidationError(c, "You do not have enough available slots to purchase this card.", nil)
	}

	tax := int64(float64(listing.Price) * 0.135)
	sellerEarnings := listing.Price - tax
	if sellerEarnings < 0 {
		sellerEarnings = 1
	}

	err = ent.WithTx(c.Request().Context(), h.app.Db(), func(tx *ent.Tx) error {
		err := tx.Player.UpdateOne(seller).
			AddLiens(sellerEarnings).
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
	if err != nil {
		return HandleErrorJson(c, err, "An error occurred while purchasing the listing.")
	}

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
				SetDescriptionf("You have sold `%s` for **%d** Liens (tax: **%d** Liens).", cardDescription, sellerEarnings, tax).
				SetEmbedFooter(h.app.Footer()).
				Build(),
		).
		Build(),
	)

	return api.SendOK(c, true)
}
