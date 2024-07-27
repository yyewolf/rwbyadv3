package topgg

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewTopGgHandler(app interfaces.App, g *echo.Group) {
	g.POST("/webhook", HandleTopGg(app))
}
