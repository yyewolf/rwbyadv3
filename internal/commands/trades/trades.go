package trades

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("Trades"),
		builder.WithEmoji(discord.Emoji{
			Name: "🔁",
		}),
		builder.WithCommands(
		//
		),
	)
}
