package dungeons

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(menus *builder.MenuStore, app interfaces.App) *builder.Menu {
	return menus.NewMenu(
		builder.WithMenuName("Dungeons"),
		builder.WithEmoji(discord.Emoji{
			Name: "🏰",
		}),
		builder.WithCommands(
			Command(menus, app),
		),
	)
}
