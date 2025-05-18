package discord

import (
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent/cookie"
)

type Options struct {
	DoRedirect     bool
	Redirect       string
	RedirectParams []string
}

type OptionsFunc func(opts *Options)

func WithRedirect(to string, params ...string) OptionsFunc {
	return func(opts *Options) {
		opts.DoRedirect = true
		opts.Redirect = to
		opts.RedirectParams = params
	}
}

func doRedirect(c echo.Context, options Options) {
	uri, _ := url.Parse("/auth/discord/")

	values := uri.Query()
	values.Add("w", options.Redirect)
	for _, param := range options.RedirectParams {
		values.Add(param, c.Param(param))
	}
	uri.RawQuery = values.Encode()

	if c.Request().Header.Get("HX-Request") != "true" {
		c.Response().Header().Set("Location", uri.String())
	}
	c.Response().Header().Set("HX-Redirect", uri.String())
}

func (h *DiscordAuthHandler) RequireAuth(opts ...OptionsFunc) func(next echo.HandlerFunc) echo.HandlerFunc {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the cookie from the request
			currentCookie, err := c.Cookie("session")
			if err != nil {
				logrus.WithError(err).Error("error getting session cookie")

				if options.DoRedirect {
					doRedirect(c, options)
					return c.NoContent(http.StatusSeeOther)
				}

				return ErrorPage(c, http.StatusUnauthorized)
			}

			// Get the session id
			sessionID := currentCookie.Value

			// Get the session
			session, err := h.app.Db().Cookie.Query().
				Where(cookie.ExpiresAtGT(time.Now())).
				Where(cookie.ID(sessionID)).
				WithPlayer().
				Only(c.Request().Context())
			if err != nil {
				logrus.WithError(err).Error("error finding session")

				if options.DoRedirect {
					doRedirect(c, options)
					return c.NoContent(http.StatusSeeOther)
				}

				return ErrorPage(c, http.StatusUnauthorized)
			}

			c.Set("session", session)

			return next(c)
		}
	}
}
