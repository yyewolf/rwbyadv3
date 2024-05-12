package boxes

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/boxes/open"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("Boxes"),
		builder.WithEmoji(discord.Emoji{
			Name: "📦",
		}),
		builder.WithCommands(
			open.OpenCommand(ms, app),
		),
	)
}
