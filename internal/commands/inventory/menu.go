package inventory

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/inventory/info"
	selectc "github.com/yyewolf/rwbyadv3/internal/commands/inventory/select"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(menus *builder.MenuStore, app interfaces.App) *builder.Menu {
	return menus.NewMenu(
		builder.WithMenuName("Inventory"),
		builder.WithEmoji(discord.Emoji{
			Name: "💼",
		}),
		builder.WithCommands(
			InventoryCommand(menus, app),
			selectc.SelectCommand(menus, app),
			info.InfoCommand(menus, app),
		),
	)
}
