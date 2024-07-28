package cdn

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func CachingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		return next(c)
	}
}

func NewCDNHandler(app interfaces.App, g *echo.Group) {
	g.Add(
		http.MethodGet,
		"/cards/*",
		echo.StaticDirectoryHandler(echo.MustSubFS(rwbyadv3.GetCardFS(), "cards/img"), false),
		CachingMiddleware,
	)

	g.Add(
		http.MethodGet,
		"/static/*",
		echo.StaticDirectoryHandler(echo.MustSubFS(rwbyadv3.GetStaticFS(), "static/dist"), false),
		CachingMiddleware,
	)
}
