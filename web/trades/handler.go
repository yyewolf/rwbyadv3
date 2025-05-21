package trades

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
)

func NewTradesHandler(app interfaces.App, g *echo.Group) {
	g.GET("/:tradeId", View(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "tradeId")))
}
