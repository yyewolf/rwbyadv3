package dungeons

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
)

func NewDungeonsHandler(app interfaces.App, g *echo.Group) {
	if app.Config().Mode == "dev" {
		g.Static("/:dungeonId/", "dungeons/dist")
	} else {
		g.StaticFS("/:dungeonId/", echo.MustSubFS(rwbyadv3.GetDungeonFS(), "dungeons/dist"))
	}

	g.GET("/:dungeonId", View(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "dungeonId")))
	g.GET("/api/dungeon", GetDebugMap(app))
	g.GET("/api/dungeon/:dungeonId", GetDungeon(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "dungeonId")))
	g.POST("/api/dungeon/:dungeonId/end", EndDungeon(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "dungeonId")))
}
