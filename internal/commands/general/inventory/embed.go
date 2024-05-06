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
	embed.SetDescription("To select a character, please use `/select`.")
	embed.SetColor(0x00ff00)

	var field discord.EmbedField
	field.Name = "Cards :"

	// Pagination here
	total := len(p.R.Cards)
	maxPage := int(math.Ceil(float64(total)/perPage)) - 1

	if page < 0 {
		page = maxPage
	}
	if page > maxPage {
		page = 0
	}

	top := (page + 1) * int(perPage)
	if top > len(p.R.Cards) {
		top = len(p.R.Cards)
	}

	displayedCards := p.R.Cards[page*int(perPage) : top]

	for i, c := range displayedCards {
		idx := page*int(perPage) + i + 1
		field.Value += fmt.Sprintf("`N°%d | %s`\n", idx, utils.Cards.FullString(c))
	}

	embed.AddFields(field)
	embed.SetFooterTextf("Page %d/%d", page+1, maxPage+1)

	customID := fmt.Sprintf("/inventory/%s/%d", p.ID, page)

	return embed.Build(), discord.NewActionRow(
		discord.NewButton(discord.ButtonStyleSecondary, "◀️ Prev", customID+"/"+componentActionPrev, ""),
		discord.NewButton(discord.ButtonStyleSecondary, "🔄 Refresh", customID+"/"+componentActionRefresh, ""),
		discord.NewButton(discord.ButtonStyleSecondary, "▶️ Next", customID+"/"+componentActionNext, ""),
	)
}
