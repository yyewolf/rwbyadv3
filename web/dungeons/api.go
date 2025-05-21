package dungeons

import (
	"context"
	"fmt"
	"math/rand"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/dungeons"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

func GetDungeon(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := utils.GetSessionFromContext(c)
		player := session.Edges.Player

		dungeonId, err := uuid.Parse(c.Param("dungeonId"))
		if err != nil {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeonId).Error(err)
			return c.JSON(500, err)
		}

		dungeon, err := app.Db().Dungeon.Get(context.Background(), dungeonId)
		if err != nil {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeonId).Error(err)
			return c.JSON(500, err)
		}

		if dungeon.PlayerID != player.ID {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeon.ID).Error("Dungeon is not owned by player")
			return c.JSON(500, err)
		}

		r := rand.New(rand.NewSource(dungeon.Seed))
		d := dungeons.NewDungeon(r)

		logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeon.ID).Info("Get Dungeon")

		return c.JSON(200, d)
	}
}

type EndDungeonRequest struct {
	Loots []string `json:"loots"`
}

func EndDungeon(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse body
		var req EndDungeonRequest
		err := c.Bind(&req)
		if err != nil {
			return c.JSON(400, err)
		}

		session := utils.GetSessionFromContext(c)
		player := session.Edges.Player

		dungeonId, err := uuid.Parse(c.Param("dungeonId"))
		if err != nil {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeonId).Error(err)
			return c.JSON(500, err)
		}

		dungeon, err := app.Db().Dungeon.Get(context.Background(), dungeonId)
		if err != nil {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeonId).Error(err)
			return c.JSON(500, err)
		}

		if dungeon.PlayerID != player.ID {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeon.ID).Error("Dungeon is not owned by player")
			return c.JSON(500, err)
		}

		r := rand.New(rand.NewSource(dungeon.Seed))
		d := dungeons.NewDungeon(r)

		var pickedUpLoots []interface{}
		err = ent.WithTx(c.Request().Context(), app.Db(), func(tx *ent.Tx) error {
			for _, loot := range d.Loots {
				if !slices.Contains(req.Loots, loot.GetID()) {
					continue
				}
				if loot.GetType() == "exit" && loot.GetID() != req.Loots[len(req.Loots)-1] {
					return fmt.Errorf("exit loot must be the last loot")
				}

				err = loot.PickedUp(tx, player)
				if err != nil {
					return err
				}

				pickedUpLoots = append(pickedUpLoots, loot)
			}

			err = tx.Dungeon.DeleteOne(dungeon).Exec(c.Request().Context())
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeon.ID).Error(err)
			return c.JSON(500, err)
		}

		texts := make([]string, 0)
		for _, loot := range loots.DungeonLoots {
			texts = append(texts, loot.RewardText(pickedUpLoots))
		}

		// remove empty texts
		texts = slices.DeleteFunc(texts, func(s string) bool {
			return s == ""
		})

		notifications.DispatchDm(app, player, discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Dungeon End").
					SetColor(app.Config().App.BotColor).
					SetDescriptionf(
						utils.Joinln(
							"Congratulations! You have completed the dungeon. Here are your rewards:",
							"",
							utils.Joinln(texts...),
						),
					).
					SetEmbedFooter(app.Footer()).
					Build(),
			).
			Build(),
		)

		logrus.WithField("user_id", player.ID).WithField("dungeon_id", dungeon.ID).WithField("loot", pickedUpLoots).Info("End Dungeon Correctly")

		return c.JSON(200, d)
	}
}

func GetDebugMap(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := rand.New(rand.NewSource(5))
		d := dungeons.NewDungeon(r)

		return c.JSON(200, d)
	}
}
