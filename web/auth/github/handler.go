package github

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/githubstar"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/errors"
	"github.com/yyewolf/rwbyadv3/web/templates/success"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubAuthHandler struct {
	app interfaces.App
	cfg *env.Config
	c   *oauth2.Config

	*echo.Group
}

func NewGithubAuthHandler(app interfaces.App, g *echo.Group) {
	cfg := app.Config()

	redirectUri, err := url.JoinPath(cfg.Github.App.BaseURI, "/callback")
	if err != nil {
		logrus.WithError(err).Fatal("error when building redirect uri")
	}

	h := &GithubAuthHandler{
		app: app,
		cfg: cfg,
		c: &oauth2.Config{
			RedirectURL:  redirectUri,
			ClientID:     cfg.Github.App.ClientID,
			ClientSecret: cfg.Github.App.ClientSecret,
			Scopes:       []string{"read:user"},
			Endpoint:     github.Endpoint,
		},
		Group: g,
	}

	g.GET("/", h.BeginAuth())
	g.GET("/callback", h.Callback())
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

func SuccessPage(c echo.Context, text string) error {
	return templates.RenderView(c, success.SuccessIndex(
		"- Auth Success",
		"",
		true,
		true,
		success.Success(text, "", ""),
	))
}

func (h *GithubAuthHandler) BeginAuth() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get state from query
		s := c.QueryParam("s")

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

		state := ImproveState(s)
		c.SetCookie(&http.Cookie{
			Name:    "oauthstate",
			Value:   state,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})
		u := h.c.AuthCodeURL(state)

		return c.Redirect(http.StatusTemporaryRedirect, u)
	}
}

func (h *GithubAuthHandler) Callback() echo.HandlerFunc {
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
		case enums.GithubCheckStar:
			return h.CallbackCheckStars(state, token)(c)
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func (h *GithubAuthHandler) CallbackCheckStars(state *ent.AuthState, token *oauth2.Token) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user
		githubUser, err := h.app.Github().GetTokenUser(token.AccessToken)
		if err != nil {
			logrus.WithField("state", state.ID).WithError(err).Error("error getting user")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		starred, err := h.app.Github().CheckTokenUserStar(token.AccessToken, h.cfg.Github.Username, h.cfg.Github.Repository)
		if err != nil {
			logrus.WithField("state", state.ID).WithError(err).Error("error checking star")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		// Send a DM
		channel, err := h.app.Client().Rest().CreateDMChannel(snowflake.MustParse(state.PlayerID))
		if err != nil {
			logrus.WithField("state", state.ID).WithError(err).Error("can't create DM channel")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		if !starred {
			h.app.Client().Rest().CreateMessage(channel.ID(),
				discord.NewMessageCreateBuilder().
					SetContentf("You have not starred the repository, you can try again by using the same command.").
					Build(),
			)

			logrus.WithField("state", state.ID).Debug("user has not starred")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		// ok, we save the user and has_starred
		err = h.app.Db().GithubStar.Update().
			Where(githubstar.PlayerID(state.PlayerID)).
			SetGithubUserID(fmt.Sprintf("%d", githubUser.GetID())).
			SetHasStarred(true).
			Exec(c.Request().Context())
		if err != nil {
			logrus.WithField("state", state.ID).WithError(err).Error("can't get github star")
			return ErrorPage(c, http.StatusInternalServerError)
		}

		h.app.Client().Rest().CreateMessage(channel.ID(),
			discord.NewMessageCreateBuilder().
				SetContentf("Your star has been taken into account!").
				Build(),
		)

		return SuccessPage(c, "Successfully looked through github, thank you, you can go back to discord now 😊")
	}
}
