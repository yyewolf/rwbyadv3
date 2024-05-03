package web

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/auth"
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
}

func (w *WebApp) Start() error {
	return w.Echo.Start(":" + w.config.Web.Port)
}

func (w *WebApp) Stop() error {
	return w.Echo.Close()
}
