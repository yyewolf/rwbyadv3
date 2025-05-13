package trades

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/dungeons"
)

func View(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		//session := utils.GetSessionFromContext(c)
		//player := session.R.Player
		tradeId := c.Param("tradeId")

		return templates.RenderView(c, dungeons.View(tradeId))
	}
}
