package market

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/market/api"
	"github.com/yyewolf/rwbyadv3/web/market/views"
)

func NewMarketHandler(app interfaces.App, g *echo.Group) {
	views.RegisterViewsRoutes(app, g)

	apiGroup := g.Group("/api")
	api.RegisterAPIRoutes(app, apiGroup)
}
