package market

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/market/listings"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(ms *builder.MenuStore, app interfaces.App) *builder.Menu {
	return ms.NewMenu(
		builder.WithMenuName("Market"),
		builder.WithEmoji(discord.Emoji{
			Name: "🛒",
		}),
		builder.WithCommands(
			listings.ListingsCommand(ms, app),
		),
	)
}
