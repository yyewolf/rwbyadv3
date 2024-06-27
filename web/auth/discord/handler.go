package discord

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/internal/values"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/errors"
	"github.com/yyewolf/rwbyadv3/web/templates/success"
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

func ErrorPage(c echo.Context, code int) error {
	return templates.RenderView(c, errors.ErrorIndex(
		"- Auth Error",
		"",
		true,
		true,
		errors.Error(fmt.Sprint(code), "Try again in a few seconds...", ""),
	))
}
func SuccessPageRedirect(c echo.Context, text, redirectURI string) error {
	return templates.RenderView(c, success.SuccessIndex(
		"- Auth Success",
		"",
		true,
		false,
		success.Success(text, "", redirectURI),
	))
}

func (h *DiscordAuthHandler) BeginAuth() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get state from query
		s := c.QueryParam("s")
		w := c.QueryParam("w")

		// w corresponds to the component calling, this tells us to create a state and redirect to the component
		if s == "" {
			state := &models.AuthDiscordState{
				State:     uuid.NewString(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				Type:      models.AuthDiscordStatesTypeLogin,
			}
			s = state.State
			switch w {
			case "main":
				state.RedirectURI = "/"
			case "market":
				state.RedirectURI = "/market"
			default:
				return ErrorPage(c, http.StatusForbidden)
			}
			err := state.InsertG(context.Background(), boil.Infer())
			if err != nil {
				logrus.WithError(err).Error("error inserting state")
				return ErrorPage(c, http.StatusInternalServerError)
			}
		} else {
			_, err := models.FindAuthDiscordStateG(context.Background(), s)
			if err != nil {
				logrus.WithField("state", s).Debug("state not found in DB")
				return ErrorPage(c, http.StatusForbidden)
			}
		}

		// Check if user is already logged in
		sessionID, err := c.Cookie("session")
		if err == nil {
			session, err := models.FindAuthCookieG(context.Background(), sessionID.Value)
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
		state, err := models.FindAuthDiscordStateG(context.Background(), s)
		if err != nil {
			logrus.WithField("state", s).Debug("state not found in DB")
			return ErrorPage(c, http.StatusForbidden)
		}

		state.DeleteG(context.Background(), false)

		token, err := h.c.Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			logrus.WithError(err).Debug("error invalid code")
			return ErrorPage(c, http.StatusForbidden)
		}

		switch state.Type {
		case models.AuthDiscordStatesTypeLogin:
			return h.CallbackLogin(state, token)(c)
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func (h *DiscordAuthHandler) CallbackLogin(state *models.AuthDiscordState, token *oauth2.Token) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user
		user, err := h.app.Client().Rest().GetCurrentUser(token.AccessToken)
		if err != nil {
			logrus.WithField("state", state.State).WithError(err).Error("error getting user")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		sessionID := utils.GenerateNewCookieId()

		// Create a session
		session := &models.AuthCookie{
			ID:        sessionID,
			PlayerID:  user.ID.String(),
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		}

		err = session.InsertG(context.Background(), boil.Infer())
		if err != nil {
			logrus.WithField("state", state.State).WithError(err).Error("error inserting session")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		cookie := &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			MaxAge:   7 * 24 * 60 * 60,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}

		// not secure on preprod
		if h.app.Config().Mode == values.Dev {
			cookie.Secure = false
		}

		// Set the cookie
		c.SetCookie(cookie)

		return SuccessPageRedirect(c, "Successfully logged in, you should be redirected in a few seconds...", state.RedirectURI)
	}
}
