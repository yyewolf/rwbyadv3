package inventory

import (
	"context"
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
)

var (
	perPage = 10.0
)

func (cmd *inventoryCommand) generator(username string, p *ent.Player, page int) (discord.Embed, discord.ContainerComponent, error) {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's inventory :", username)
	embed.SetDescriptionf("To select a character, please use %s.", cmd.app.CommandMention("select"))
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	count, err := p.QueryCards().
		Where(card.Available(true)).
		Count(context.Background())
	if err != nil {
		return discord.Embed{}, nil, err
	}

	total := count
	maxPage := int(math.Ceil(float64(total)/perPage)) - 1

	var field discord.EmbedField

	if page < 0 {
		page = maxPage
	}
	if page > maxPage {
		page = 0
	}

	top := (page + 1) * int(perPage)
	if top > count {
		top = count
	}

	cards, err := p.QueryCards().
		Order(card.ByPosition()).
		Where(card.Available(true)).
		Offset(page * int(perPage)).
		Limit(int(perPage)).
		All(context.Background())
	if err != nil {
		return discord.Embed{}, nil, err
	}

	field.Name = fmt.Sprintf("Cards (page %d/%d) :", page+1, maxPage+1)

	for i, c := range cards {
		idx := page*int(perPage) + i + 1
		field.Value += fmt.Sprintf("`N°%d | %s`\n", idx, c.FullString())
	}

	if len(cards) == 0 {
		field.Name = "Cards :"
		field.Value = "You have no cards to be shown."
	}

	embed.AddFields(field)

	customID := fmt.Sprintf("/inventory/%s/%d", p.ID, page)

	return embed.Build(), discord.NewActionRow(
		discord.NewSecondaryButton("◀️ Prev", customID+"/"+componentActionPrev),
		discord.NewSecondaryButton("🔄 Refresh", customID+"/"+componentActionRefresh),
		discord.NewSecondaryButton("▶️ Next", customID+"/"+componentActionNext),
	), nil
}
