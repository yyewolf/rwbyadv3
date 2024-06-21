package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
	"github.com/yyewolf/rwbyadv3/web/auth/github"
)

var (
	DiscordHandler *discord.DiscordAuthHandler
)

func NewAuthHandler(app interfaces.App, g *echo.Group) {
	github.NewGithubAuthHandler(app, g.Group("/github"))
	DiscordHandler = discord.NewDiscordAuthHandler(app, g.Group("/discord"))

	g.GET("/discord", func(c echo.Context) error {
		return c.Redirect(302, "/discord/")
	})
	g.GET("/github", func(c echo.Context) error {
		return c.Redirect(302, "/github/")
	})
}
