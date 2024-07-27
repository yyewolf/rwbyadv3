package dungeons

import (
	"context"
	"math/rand"
	"slices"

	"github.com/disgoorg/disgo/discord"
	"github.com/labstack/echo/v4"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/dungeons"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func GetDungeon(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := utils.GetSessionFromContext(c)
		player := session.R.Player

		dungeon, err := models.FindDungeonG(context.Background(), c.Param("dungeonId"))
		if err != nil {
			return c.JSON(500, err)
		}

		if dungeon.PlayerID != player.ID {
			return c.JSON(500, err)
		}

		r := rand.New(rand.NewSource(dungeon.Seed))
		d := dungeons.NewDungeon(r)

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
		player := session.R.Player

		dungeon, err := models.FindDungeonG(context.Background(), c.Param("dungeonId"))
		if err != nil {
			return c.JSON(500, err)
		}

		if dungeon.PlayerID != player.ID {
			return c.JSON(500, err)
		}

		r := rand.New(rand.NewSource(dungeon.Seed))
		d := dungeons.NewDungeon(r)

		tx, err := boil.BeginTx(context.Background(), nil)
		if err != nil {
			return c.JSON(500, err)
		}

		var pickedUpLoots []interface{}

		for _, loot := range d.Loots {
			if !slices.Contains(req.Loots, loot.GetID()) {
				continue
			}
			if loot.GetType() == "exit" && loot.GetID() != req.Loots[len(req.Loots)-1] {
				tx.Rollback()
				return c.JSON(500, err)
			}
			loot.PickedUp(tx, player)
			pickedUpLoots = append(pickedUpLoots, loot)
		}

		player.Update(context.Background(), tx, boil.Infer())
		dungeon.Delete(context.Background(), tx, false)

		err = tx.Commit()
		if err != nil {
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
