package market

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/market/api"
)

func NewMarketHandler(app interfaces.App, g *echo.Group) {
	api.RegisterAPIRoutes(app, g)
}
