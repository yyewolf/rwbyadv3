package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth"
	"github.com/yyewolf/rwbyadv3/web/cdn"
	"github.com/yyewolf/rwbyadv3/web/market"
	"github.com/yyewolf/rwbyadv3/web/metrics"
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
	auth.NewAuthHandler(w.app, w.Group("/auth"))
	metrics.NewMetricsHandler(w.app, w.Group("/metrics"))

	// Also redirect from /market to /market/
	w.GET("/market", RedirectTo("/market/"))
	market.NewMarketHandler(w.app, w.Group("/market"))

	cdn.NewCDNHandler(w.app, w.Group("/cdn"))
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
