package discord

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/yyewolf/rwbyadv3/models"
)

type Options struct {
	DoRedirect bool
	Redirect   string
}

type OptionsFunc func(opts *Options)

func WithRedirect(to string) OptionsFunc {
	return func(opts *Options) {
		opts.DoRedirect = true
		opts.Redirect = to
	}
}

func (h *DiscordAuthHandler) RequireAuth(opts ...OptionsFunc) func(next echo.HandlerFunc) echo.HandlerFunc {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the cookie from the request
			cookie, err := c.Cookie("session")
			if err != nil {
				logrus.WithError(err).Error("error getting session cookie")

				if options.DoRedirect {
					if c.Request().Header.Get("HX-Request") != "true" {
						c.Response().Header().Set("Location", fmt.Sprintf("/auth/discord/?w=%s", options.Redirect))
					}
					c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/auth/discord/?w=%s", options.Redirect))
					return c.NoContent(http.StatusSeeOther)
				}

				return ErrorPage(c, http.StatusUnauthorized)
			}

			// Get the session id
			sessionID := cookie.Value

			// Get the session
			session, err := models.AuthCookies(
				qm.Where("expires_at > NOW()"),
				qm.Where(models.AuthCookieColumns.ID+"=?", sessionID),
				qm.Load(models.AuthCookieRels.Player),
			).OneG(context.Background())
			if err != nil {
				logrus.WithError(err).Error("error finding session")

				if options.DoRedirect {
					if c.Request().Header.Get("HX-Request") != "true" {
						c.Response().Header().Set("Location", fmt.Sprintf("/auth/discord/?w=%s", options.Redirect))
					}
					c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/auth/discord/?w=%s", options.Redirect))
					return c.NoContent(http.StatusSeeOther)
				}

				return ErrorPage(c, http.StatusUnauthorized)
			}

			c.Set("session", session)

			return next(c)
		}
	}
}
