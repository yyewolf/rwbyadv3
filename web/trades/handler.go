package trades

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/api"
)

type TradeApiHandler struct {
	app interfaces.App
}

func HandleErrorJson(c echo.Context, err error, userMessage string) error {
	logrus.WithError(err).Error(userMessage)
	return api.NewErrorResponse[any](api.ErrorInternalServerError, userMessage).JSON(c, http.StatusInternalServerError)
}

func NewTradesHandler(app interfaces.App, g *echo.Group) {
	// handler := TradeApiHandler{app: app}

	// g.GET("/api/my-cards", api.MyCards(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectTrades, "tradeId")))
	// g.GET("/:tradeId/api/their-cards", api.TheirCards(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectTrades, "tradeId")))
}
