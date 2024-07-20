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
	Loots []string `json:"string"`
}

func EndDungeon(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse body
		var req EndDungeonRequest
		c.Bind(&req)

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

		for _, loot := range d.Loots {
			if !slices.Contains(req.Loots, loot.GetID()) {
				continue
			}
			loot.PickedUp(tx, player)
		}

		player.Update(context.Background(), tx, boil.Infer())
		dungeon.Delete(context.Background(), tx, false)

		err = tx.Commit()
		if err != nil {
			return c.JSON(500, err)
		}

		notifications.DispatchDm(app, player, discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitle("Dungeon End").
					SetColor(app.Config().App.BotColor).
					SetDescriptionf("Yes").
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
