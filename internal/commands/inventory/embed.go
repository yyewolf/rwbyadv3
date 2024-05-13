package inventory

import (
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

var (
	perPage = 20.0
)

func (cmd *inventoryCommand) generator(username string, p *models.Player, page int) (discord.Embed, discord.ContainerComponent) {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's inventory :", username)
	embed.SetDescriptionf("To select a character, please use %s.", cmd.app.CommandMention("select"))
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	total := len(p.R.PlayerCards)
	maxPage := int(math.Ceil(float64(total)/perPage)) - 1

	var field discord.EmbedField
	field.Name = fmt.Sprintf("Cards (page %d/%d) :", page+1, maxPage+1)

	if page < 0 {
		page = maxPage
	}
	if page > maxPage {
		page = 0
	}

	top := (page + 1) * int(perPage)
	if top > len(p.R.PlayerCards) {
		top = len(p.R.PlayerCards)
	}

	displayedCards := p.R.PlayerCards[page*int(perPage) : top]

	for i, pre := range displayedCards {
		c := pre.R.Card
		idx := page*int(perPage) + i + 1
		field.Value += fmt.Sprintf("`N°%d | %s`\n", idx, utils.Cards.FullString(c))
	}

	if len(displayedCards) == 0 {
		field.Name = "Cards :"
		field.Value = "You have no cards to be shown."
	}

	embed.AddFields(field)

	customID := fmt.Sprintf("/inventory/%s/%d", p.ID, page)

	return embed.Build(), discord.NewActionRow(
		discord.NewButton(discord.ButtonStyleSecondary, "◀️ Prev", customID+"/"+componentActionPrev, ""),
		discord.NewButton(discord.ButtonStyleSecondary, "🔄 Refresh", customID+"/"+componentActionRefresh, ""),
		discord.NewButton(discord.ButtonStyleSecondary, "▶️ Next", customID+"/"+componentActionNext, ""),
	)
}
