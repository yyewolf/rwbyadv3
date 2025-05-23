package web

import (
	"net/http"
	"os"

	"github.com/yyewolf/rwbyadv3"
	"github.com/yyewolf/rwbyadv3/web/trades"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/cdn"
	"github.com/yyewolf/rwbyadv3/web/dungeons"
	"github.com/yyewolf/rwbyadv3/web/landing"
	"github.com/yyewolf/rwbyadv3/web/market"
	"github.com/yyewolf/rwbyadv3/web/metrics"
	"github.com/yyewolf/rwbyadv3/web/topgg"
)

type WebApp struct {
	app    interfaces.App
	config *env.Config

	*echo.Echo
}

func NewWebApp(opts ...Option) *WebApp {
	var app = &WebApp{}

	for _, opt := range opts {
		opt(app)
	}

	if app.Echo == nil {
		app.Echo = echo.New()
	}

	app.RegisterRoutes()

	return app
}

func (w *WebApp) RegisterRoutes() {
	apis := w.Group("/apis")

	auth.NewAuthHandler(w.app, w.Group("/auth"))
	metrics.NewMetricsHandler(w.app, apis.Group("/metrics"))
	market.NewMarketHandler(w.app, apis.Group("/market"))
	dungeons.NewDungeonsHandler(w.app, w.Group("/dungeons"))
	trades.NewTradesHandler(w.app, apis.Group("/trades"))

	cdn.NewCDNHandler(w.app, w.Group("/cdn"))

	// Register landing error handler
	landing.RegisterErrorHandler(w.app, w.Group("/landing"))

	topgg.NewTopGgHandler(w.app, w.Group("/topgg"))

	// If No Route, go to the static files :
	fs := echo.MustSubFS(rwbyadv3.GetWwwFS(), "www/build")
	if w.app.Config().Mode == "dev" {
		fs = echo.MustSubFS(os.DirFS("."), "www/build-dev")
	}

	w.RouteNotFound("*", echo.StaticDirectoryHandler(fs, false))
}

func (w *WebApp) Start() error {
	return w.Echo.Start(":" + w.config.Web.Port)
}

func (w *WebApp) Stop() error {
	return w.Echo.Close()
}

func RedirectTo(path string) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, path)
	}
}
