package inventory

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	selectc "github.com/yyewolf/rwbyadv3/internal/commands/inventory/select"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("Inventory"),
		builder.WithEmoji(discord.Emoji{
			Name: "💼",
		}),
		builder.WithCommands(
			InventoryCommand(ms, app),
			selectc.SelectCommand(ms, app),
		),
	)
}
