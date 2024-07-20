package dungeons

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/web/templates"
	"github.com/yyewolf/rwbyadv3/web/templates/dungeons"
	"github.com/yyewolf/rwbyadv3/web/templates/errors"
)

func View(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := utils.GetSessionFromContext(c)
		player := session.R.Player
		dungeonId := c.Param("dungeonId")

		dungeon, err := models.FindDungeonG(context.Background(), dungeonId)
		if err != nil {
			return templates.RenderView(c, errors.ErrorIndex(
				"- Dungeon Error",
				"",
				true,
				true,
				errors.Error("404", "Dungeon not found !", ""),
			))
		}

		if dungeon.PlayerID != player.ID {
			return templates.RenderView(c, errors.ErrorIndex(
				"- Dungeon Error",
				"",
				true,
				true,
				errors.Error("404", "Dungeon not found !", ""),
			))
		}

		return templates.RenderView(c, dungeons.View(dungeonId))
	}
}
