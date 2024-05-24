package begin

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func (cmd *beginCommand) generator(p *models.Player, page int) (discord.Embed, discord.ContainerComponent) {
	embed := discord.NewEmbedBuilder()
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	maxPage := 2

	if page < 0 {
		page = maxPage
	}
	if page > maxPage {
		page = 0
	}

	embed.SetTitlef("Welcome to RWBY Adventures (%d/%d)", page+1, maxPage+1)

	switch page {
	case 0:
		cmd.pageOne(embed)
	case 1:
		cmd.pageTwo(embed)
	case 2:
		cmd.pageThree(embed)
	}

	customID := fmt.Sprintf("/begin/%s/%d", p.ID, page)

	return embed.Build(), discord.NewActionRow(
		discord.NewButton(discord.ButtonStyleSecondary, "◀️ Prev", customID+"/"+componentActionPrev, ""),
		discord.NewButton(discord.ButtonStyleSecondary, "▶️ Next", customID+"/"+componentActionNext, ""),
	)
}

func (cmd *beginCommand) pageOne(e *discord.EmbedBuilder) {
	e.AddField(
		"What is this?",
		utils.Joinln(
			"RWBY Adventures is a text-based adventure game set in the world of RWBY.",
			"Players can collect *cards*, battle each others in *duels*, take on dungeons, and more!",
			fmt.Sprintf("Using the %s command, you can always bring up a useful tooltip.", cmd.app.CommandMention("help")),
		),
		false,
	)

	e.AddField(
		"Open-Source",
		utils.Joinln(
			"RWBY Adventures is also [open-source](https://github/yyewolf/rwbyadv3).",
			"A [documentation]() is also available for technical help and command insights.",
			"The documentation can be used to go further in-depth with the game mechanics.",
		),
		false,
	)

	// Disclaimer
	e.AddField(
		"Disclaimer",
		"*This is not endorsed by Rooster Teeth in any way. Views, opinions, and thoughts are all my own. Rooster Teeth and RWBY are trade names or registered trademarks of Rooster Teeth Productions, LLC. © Rooster Teeth Productions, LLC.*",
		false,
	)
}

func (cmd *beginCommand) pageTwo(e *discord.EmbedBuilder) {
	e.AddField(
		"Cards",
		utils.Joinln(
			"Not all of the characters from RWBY are implemented yet, but as time goes on, new cards will be added!",
			"Cards can also be grimms or other entities from the RWBY universe.",
		),
		false,
	)

	e.AddField(
		"Loot boxes",
		utils.Joinln(
			"Naturally, you can open loot boxes to obtain cards.",
			fmt.Sprintf("You can open a loot box using the %s command.", cmd.app.CommandMention("open")),
			"You can even try it now!",
		),
		false,
	)
}

func (cmd *beginCommand) pageThree(e *discord.EmbedBuilder) {
	e.AddField(
		"Battles",
		utils.Joinln(
			"Players can battle each other in duels.",
			"Each player has a deck of cards, and they can use them to fight.",
			"Not all cards can be used in battles, this is to avoid issues with non-fighting characters.",
			"Zwei on the other hand, is combat ready!",
		),
		false,
	)

	e.AddField(
		"Trading",
		utils.Joinln(
			"Of course, players can trade cards with each other.",
			"There's even a marketplace where players can sell their cards to others.",
			"Also, in the marketplace, you can put cards up for auction!",
		),
		false,
	)

	e.AddField(
		"The end",
		utils.Joinln(
			"This does not represent all of the features, but rather a quick getting started!",
			"Feel free to explore the game and have fun!",
			"The support server is also available for any questions or feedback.",
		),
		false,
	)
}
