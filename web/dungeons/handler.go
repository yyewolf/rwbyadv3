package dungeons

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/maze"
)

func GetMapHandler(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		m := maze.DefaultMaze(15, 15)
		m.Generate()

		return c.JSON(200, m.Expand(3))
	}
}

func NewDungeonsHandler(app interfaces.App, g *echo.Group) {
	if app.Config().Mode == "dev" {
		g.Static("/", "dungeons/dist")
	} else {
		g.StaticFS("/", echo.MustSubFS(rwbyadv3.GetDungeonFS(), "dungeons/dist"))
	}

	g.GET("/map", GetMapHandler(app))
}
