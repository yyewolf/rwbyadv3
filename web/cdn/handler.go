package cdn

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewCDNHandler(app interfaces.App, g *echo.Group) {
	g.StaticFS("/cards", echo.MustSubFS(rwbyadv3.GetCardFS(), "cards/img"))
}
