package trades

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/api"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
)

type TradeApiHandler struct {
	app interfaces.App
}

func HandleErrorJson(c echo.Context, err error, userMessage string) error {
	logrus.WithError(err).Error(userMessage)
	return api.NewErrorResponse[any](api.ErrorInternalServerError, userMessage).JSON(c, http.StatusInternalServerError)
}

func NewTradesHandler(app interfaces.App, g *echo.Group) {
	handler := TradeApiHandler{app: app}

	g.GET("/self/cards", handler.MyCards(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectTrades)))
	g.POST("/:playerId", handler.CreateTrade(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectTrades, "playerId")))
	g.GET("/:playerId/cards", handler.TheirCards(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectTrades, "playerId")))
}
