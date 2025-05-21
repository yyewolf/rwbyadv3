package market

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/commands/market/auctions"
	"github.com/yyewolf/rwbyadv3/internal/commands/market/listings"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

func NewMenu(menus *builder.MenuStore, app interfaces.App) *builder.Menu {
	return menus.NewMenu(
		builder.WithMenuName("Market"),
		builder.WithEmoji(discord.Emoji{
			Name: "🛒",
		}),
		builder.WithCommands(
			listings.ListingsCommand(menus, app),
			auctions.AuctionsCommand(menus, app),
		),
	)
}
