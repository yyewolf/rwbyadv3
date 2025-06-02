package trades

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/api"
	"github.com/yyewolf/rwbyadv3/web/landing"
)

func (h *TradeApiHandler) TheirCards(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		player, err := app.Db().Player.Get(c.Request().Context(), c.Param("playerId"))
		if err != nil {
			if ent.IsNotFound(err) {
				return landing.APIErrorPage(c, landing.ErrorNotFound)
			}
			return landing.APIErrorPage(c, landing.ErrorTrade)
		}

		playerCards, err := player.QueryCards().
			Where(
				card.Available(true),
			).
			WithType().
			Order(card.ByPosition()).
			All(c.Request().Context())
		if err != nil {
			return landing.APIErrorPage(c, landing.ErrorTrade)
		}

		return api.SendOK(c, ent.ViewCardListAs(playerCards, ent.Public))
	}
}
