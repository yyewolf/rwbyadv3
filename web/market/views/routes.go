package views

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/templates/market"
)

func RegisterViewsRoutes(app interfaces.App, g *echo.Group) {
	g.GET("/", echo.WrapHandler(templ.Handler(market.Main())))
}
