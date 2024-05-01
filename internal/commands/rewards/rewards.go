package rewards

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/rewards/stars"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("Rewards"),
		builder.WithEmoji(discord.Emoji{
			Name: "🎁",
		}),
		builder.WithCommands(
			stars.StarCommand(ms, app),
		),
	)
}
