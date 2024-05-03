package bugs

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/bugs/report"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("Bugs"),
		builder.WithEmoji(discord.Emoji{
			Name: "🐞",
		}),
		builder.WithCommands(
			report.ReportCommand(ms, app),
		),
	)
}
