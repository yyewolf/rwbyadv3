package inventory

import (
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

var (
	perPage = 10.0
)

func (cmd *inventoryCommand) generator(username string, p *models.Player, page int) (discord.Embed, discord.ContainerComponent) {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's inventory :", username)
	embed.SetDescriptionf("To select a character, please use %s.", cmd.app.CommandMention("select"))
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	cards := utils.Players.AvailableCards(p)
	total := len(cards)
	maxPage := int(math.Ceil(float64(total)/perPage)) - 1

	var field discord.EmbedField

	if page < 0 {
		page = maxPage
	}
	if page > maxPage {
		page = 0
	}

	top := (page + 1) * int(perPage)
	if top > len(cards) {
		top = len(cards)
	}

	field.Name = fmt.Sprintf("Cards (page %d/%d) :", page+1, maxPage+1)

	displayedCards := cards[page*int(perPage) : top]

	for i, c := range displayedCards {
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
		discord.NewSecondaryButton("◀️ Prev", customID+"/"+componentActionPrev),
		discord.NewSecondaryButton("🔄 Refresh", customID+"/"+componentActionRefresh),
		discord.NewSecondaryButton("▶️ Next", customID+"/"+componentActionNext),
	)
}
