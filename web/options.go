package web

import (
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type Option func(a *WebApp)

func WithConfig(config *env.Config) Option {
	return func(a *WebApp) {
		a.config = config
	}
}

func WithApp(app interfaces.App) Option {
	return func(a *WebApp) {
		a.app = app
		a.config = app.Config()
	}
}

func WithEcho(e *echo.Echo) Option {
	return func(a *WebApp) {
		a.Echo = e
	}
}
