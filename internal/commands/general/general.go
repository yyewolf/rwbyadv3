package general

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/general/begin"
	"github.com/yyewolf/rwbyadv3/internal/commands/general/help"
	"github.com/yyewolf/rwbyadv3/internal/commands/general/profile"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(menus *builder.MenuStore, app interfaces.App) *builder.Menu {
	return menus.NewMenu(
		builder.WithMenuName("General"),
		builder.WithEmoji(discord.Emoji{
			Name: "💻",
		}),
		builder.WithCommands(
			help.HelpCommand(menus, app),
			profile.ProfileCommand(menus, app),
			begin.BeginCommand(menus, app),
		),
	)
}
