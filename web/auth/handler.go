package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
	"github.com/yyewolf/rwbyadv3/web/auth/github"
)

func NewAuthHandler(app interfaces.App, g *echo.Group) {
	github.NewGithubAuthHandler(app, g.Group("/github"))
	discord.NewDiscordAuthHandler(app, g.Group("/discord"))
}
