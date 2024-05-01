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
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
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

func (h *GithubAuthHandler) BeginAuth() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get state from query
		s := c.QueryParam("s")

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
		oauthState, _ := c.Cookie("oauthstate")
		if c.FormValue("state") != oauthState.Value {
			logrus.Warn("differing state")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		s := ReverseState(oauthState.Value)
		state, err := models.FindAuthGithubStateG(context.Background(), s)
		if err != nil {
			logrus.WithField("state", s).Error("state not found in DB")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		state.DeleteG(context.Background(), false)

		token, err := h.c.Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			logrus.WithError(err).Error("error invalid code")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		switch state.Type {
		case models.AuthGithubStatesTypeCheckStar:
			return h.CallbackCheckStars(state, token)(c)
		}
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func (h *GithubAuthHandler) CallbackCheckStars(state *models.AuthGithubState, token *oauth2.Token) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user
		githubUser, err := h.app.Github().GetTokenUser(token.AccessToken)
		if err != nil {
			logrus.WithField("state", state.State).WithError(err).Error("error getting user")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		starred, err := h.app.Github().CheckTokenUserStar(token.AccessToken, h.cfg.Github.Username, h.cfg.Github.Repository)
		if err != nil {
			logrus.WithField("state", state.State).WithError(err).Error("error checking star")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		// Send a DM
		channel, err := h.app.Client().Rest().CreateDMChannel(snowflake.MustParse(state.PlayerID))
		if err != nil {
			logrus.WithField("state", state.State).WithError(err).Error("can't create DM channel")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		if !starred {
			h.app.Client().Rest().CreateMessage(channel.ID(),
				discord.NewMessageCreateBuilder().
					SetContentf("You have not starred the repository, you can try again by using the same command.").
					Build(),
			)

			logrus.WithField("state", state.State).Error("user has not starred")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		// ok, we save the user and has_starred
		githubStar, err := models.FindGithubStarG(context.Background(), state.PlayerID)
		if err != nil {
			logrus.WithField("state", state.State).WithError(err).Error("can't get github star")
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}
		githubStar.HasStarred = true
		githubStar.GithubUserID = null.NewString(fmt.Sprintf("%d", githubUser.GetID()), true)
		githubStar.UpdateG(context.Background(), boil.Infer())

		h.app.Client().Rest().CreateMessage(channel.ID(),
			discord.NewMessageCreateBuilder().
				SetContentf("Your star has been taken into account!").
				Build(),
		)

		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
