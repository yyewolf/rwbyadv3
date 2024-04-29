package general

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/general/help"
	"github.com/yyewolf/rwbyadv3/internal/commands/general/profile"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("General"),
		builder.WithEmoji(discord.Emoji{
			Name: "💻",
		}),
		builder.WithCommands(
			help.HelpCommand(ms, app),
			profile.ProfileCommand(ms, app),
		),
	)
}
