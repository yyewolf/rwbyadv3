package dungeons

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewDungeonsHandler(app interfaces.App, g *echo.Group) {
	g.StaticFS("/", echo.MustSubFS(rwbyadv3.GetDungeonFS(), "dungeons"))
}
