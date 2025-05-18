package topgg

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

type botVote struct {
	Bot       string `json:"bot"`
	User      string `json:"user"`
	Type      string `json:"type"`
	IsWeekend bool   `json:"isWeekend"`
	Query     string `json:"query?"`
}

func HandleTopGg(app interfaces.App) func(c echo.Context) error {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") != app.Config().TopGg.Token {
			logrus.
				WithField("from", "topgg").
				Error("Wrong token received for vote.")
			return echo.ErrForbidden
		}

		var req botVote

		err := c.Bind(&req)
		if err != nil {
			logrus.
				WithField("from", "topgg").
				WithError(err).
				Error("Can't bind request body.")
			return echo.ErrForbidden
		}

		userID, err := snowflake.Parse(req.User)
		if err != nil {
			logrus.
				WithField("from", "topgg").
				WithField("req", req).
				WithError(err).
				Error("Can't parse user id.")
			return echo.ErrForbidden
		}

		ctx, err := builder.GetContext(
			context.Background(),
			app,
			userID,
			builder.WithPlayer(),
			builder.WithPlayerDaily(),
		)
		if err != nil {
			logrus.
				WithField("from", "topgg").
				WithField("req", req).
				WithError(err).
				Error("Can't get user from ctx.")
			return c.JSON(200, "ok")
		}

		player := ctx.Value(builder.NewPlayerKey).(*ent.Player)

		// Check if we increment or reset streak
		if time.Since(player.Edges.Daily.LastVoteAt) > 24*time.Hour {
			player.Edges.Daily.Streak = 0
		}
		player.Edges.Daily.Streak++

		err = player.Edges.Daily.Update().
			SetHasVoted(true).
			SetLastVoteAt(time.Now()).
			SetStreak(player.Edges.Daily.Streak).
			Exec(ctx)

		if err != nil {
			logrus.
				WithField("from", "topgg").
				WithField("req", req).
				WithError(err).
				Error("Can't store vote into database.")
			return echo.ErrInternalServerError
		}

		logrus.
			WithField("from", "topgg").
			WithField("player", req.User).
			Info("Got vote from player.")

		notifications.DispatchDm(app, player, discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Daily Reward :").
					SetColor(app.Config().App.BotColor).
					SetDescriptionf(
						utils.Joinln(
							"Thank you for your vote !",
							fmt.Sprintf("You can claim your reward with %s !", app.CommandMention("daily")),
						),
					).
					SetEmbedFooter(app.Footer()).
					Build(),
			).
			Build(),
		)

		return c.JSON(200, "ok")
	}
}
