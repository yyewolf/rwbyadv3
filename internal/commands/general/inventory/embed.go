package inventory

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

func (cmd *inventoryCommand) generateEmbed(e *handler.CommandEvent, p *models.Player) discord.Embed {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's Inventory", e.User().Username)
	embed.SetDescription("To select a character, please use `/select`.")
	embed.SetColor(0x00ff00)

	var field discord.EmbedField
	field.Name = "Cards :"

	for i, c := range p.R.Cards {
		field.Value += fmt.Sprintf("`N°%d | %s`\n", i+1, utils.Cards.FullString(c))
	}

	embed.AddFields(field)

	return embed.Build()
}
