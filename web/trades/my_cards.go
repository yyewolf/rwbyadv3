package trades

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/api"
)

func (h *TradeApiHandler) MyCards(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := utils.GetSessionFromContext(c)
		player := session.Edges.Player

		playerCards, err := player.QueryCards().
			Where(
				card.Available(true),
			).
			WithType().
			Order(card.ByPosition()).
			All(c.Request().Context())
		if err != nil {
			return HandleErrorJson(c, err, "failed to fetch player cards")
		}

		return api.SendOK(c, ent.ViewCardListAs(playerCards, ent.Self))
	}
}
