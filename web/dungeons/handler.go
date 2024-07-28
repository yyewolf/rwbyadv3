package dungeons

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/auth/discord"
)

func CachingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		return next(c)
	}
}

func NewDungeonsHandler(app interfaces.App, g *echo.Group) {
	fs := echo.MustSubFS(rwbyadv3.GetDungeonFS(), "dungeons/dist")
	if app.Config().Mode == "dev" {
		fs = echo.MustSubFS(os.DirFS("."), "dungeons/dist")
	}

	g.Add(
		http.MethodGet,
		"/static/*",
		echo.StaticDirectoryHandler(fs, false),
		CachingMiddleware,
	)

	g.GET("/:dungeonId", View(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "dungeonId")))
	g.GET("/api/dungeon", GetDebugMap(app))
	g.GET("/api/dungeon/:dungeonId", GetDungeon(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "dungeonId")))
	g.POST("/api/dungeon/:dungeonId/end", EndDungeon(app), auth.DiscordHandler.RequireAuth(discord.WithRedirect(discord.RedirectDungeons, "dungeonId")))
}
