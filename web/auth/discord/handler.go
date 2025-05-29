package discord

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"golang.org/x/oauth2"
)

type DiscordAuthHandler struct {
	app interfaces.App
	cfg *env.Config
	c   *oauth2.Config

	*echo.Group
}

func NewDiscordAuthHandler(app interfaces.App, g *echo.Group) *DiscordAuthHandler {
	cfg := app.Config()

	redirectUri, err := url.JoinPath(cfg.Discord.App.BaseURI, "/callback")
	if err != nil {
		logrus.WithError(err).Fatal("error when building redirect uri")
	}

	h := &DiscordAuthHandler{
		app: app,
		cfg: cfg,
		c: &oauth2.Config{
			RedirectURL:  redirectUri,
			ClientID:     cfg.Discord.App.ClientID,
			ClientSecret: cfg.Discord.App.ClientSecret,
			Scopes:       []string{"identify"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		},
		Group: g,
	}

	g.GET("/", h.BeginAuth())
	g.GET("/callback", h.Callback())

	return h
}

func ImproveState(s string) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return s + "/" + state
}

func ReverseState(s string) string {
	return strings.Split(s, "/")[0]
}

func (h *DiscordAuthHandler) BeginAuth() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get state from query
		params := c.QueryParams()

		s := params.Get("s")
		w := params.Get("w")

		// w corresponds to the component calling, this tells us to create a state and redirect to the component
		if s == "" {
			redirectUri := "/"
			switch w {
			case RedirectMain:
				redirectUri = "/redirect/?to=main"
			case RedirectMarket:
				redirectUri = "/redirect/?to=market"
			case RedirectDungeons:
				redirectUri = "/redirect/?to=dungeons&id=" + params.Get("dungeonId")
			case RedirectTrades:
				redirectUri = "/redirect/?to=trades&id=" + params.Get("tradeId")
			default:
				return ErrorPage(c, http.StatusForbidden)
			}

			state, err := h.app.Db().AuthState.Create().
				SetExpiresAt(time.Now().Add(24 * time.Hour)).
				SetType(enums.DiscordLogin).
				SetRedirectURI(redirectUri).
				Save(c.Request().Context())
			if err != nil {
				logrus.WithError(err).Error("error inserting state")
				return ErrorPage(c, http.StatusInternalServerError)
			}

			s = state.ID.String()
		} else {
			stateId, err := uuid.Parse(s)
			if err != nil {
				logrus.WithError(err).Debug("error parsing state")
				return ErrorPage(c, http.StatusForbidden)
			}

			_, err = h.app.Db().AuthState.Get(c.Request().Context(), stateId)
			if err != nil {
				logrus.WithField("state", s).Debug("state not found in DB")
				return ErrorPage(c, http.StatusForbidden)
			}
		}

		// Check if user is already logged in
		sessionID, err := c.Cookie("session")
		if err == nil {
			session, err := h.app.Db().Cookie.Get(c.Request().Context(), sessionID.Value)
			if err == nil && session.ExpiresAt.After(time.Now()) {
				redirectUri, err := utils.GetRedirectForW(w)
				if err != nil {
					c.Redirect(http.StatusTemporaryRedirect, "/")
				}
				return c.Redirect(http.StatusTemporaryRedirect, redirectUri)
			}
		}

		state := ImproveState(s)
		c.SetCookie(&http.Cookie{
			Name:  "oauthstate",
			Value: state,
		})
		u := h.c.AuthCodeURL(state)

		return c.Redirect(http.StatusTemporaryRedirect, u)
	}
}

func (h *DiscordAuthHandler) Callback() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Read oauthState from Cookie
		oauthState, err := c.Cookie("oauthstate")
		if err != nil || c.FormValue("state") != oauthState.Value {
			return ErrorPage(c, http.StatusForbidden)
		}

		s := ReverseState(oauthState.Value)

		state, err := h.app.Db().AuthState.Get(c.Request().Context(), uuid.MustParse(s))
		if err != nil {
			logrus.WithField("state", s).Debug("state not found in DB")
			return ErrorPage(c, http.StatusForbidden)
		}

		h.app.Db().AuthState.DeleteOne(state).Exec(c.Request().Context())

		token, err := h.c.Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			logrus.WithError(err).Debug("error invalid code")
			return ErrorPage(c, http.StatusForbidden)
		}

		switch state.Type {
		case enums.DiscordLogin:
			return h.CallbackLogin(state, token)(c)
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func (h *DiscordAuthHandler) CallbackLogin(state *ent.AuthState, token *oauth2.Token) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user
		user, err := h.app.Client().Rest().GetCurrentUser(token.AccessToken)
		if err != nil {
			logrus.WithField("state", state.ID).WithError(err).Error("error getting user")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		sessionID := utils.GenerateNewCookieId()

		err = h.app.Db().Cookie.Create().
			SetID(sessionID).
			SetPlayerID(user.ID.String()).
			SetExpiresAt(time.Now().Add(7 * 24 * time.Hour)).
			Exec(c.Request().Context())
		if err != nil {
			logrus.WithField("state", state.ID).WithError(err).Error("error inserting session")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		cookie := &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			MaxAge:   7 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode, // Lax is important to allow auto auth when coming from discord
		}

		// Set the cookie
		c.SetCookie(cookie)

		// return SuccessPageRedirect(c, "Successfully logged in, you should be redirected in a few seconds...", state.RedirectURI)
		return c.Redirect(http.StatusTemporaryRedirect, state.RedirectURI)
	}
}
