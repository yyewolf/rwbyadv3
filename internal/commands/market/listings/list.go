package listings

import (
	"fmt"
	"math"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

var (
	perPage = 10.0
)

func (cmd *listingsCommand) generator(username string, p *models.Player, page int) (discord.Embed, discord.ContainerComponent) {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's listings :", username)
	// embed.SetDescriptionf("To select a character, please use %s.", cmd.app.CommandMention("select"))
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	cards := utils.Players.MarketListings(p)
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

	customID := fmt.Sprintf("/listings/%s/%d", p.ID, page)

	return embed.Build(), discord.NewActionRow(
		discord.NewSecondaryButton("◀️ Prev", customID+"/"+componentActionPrev),
		discord.NewSecondaryButton("🔄 Refresh", customID+"/"+componentActionRefresh),
		discord.NewSecondaryButton("▶️ Next", customID+"/"+componentActionNext),
	)
}

func (cmd *listingsCommand) GetListings(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	username := e.User().Username
	if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	}

	embed, components := cmd.generator(username, p, 0)

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *listingsCommand) HandleGetListingsInteraction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	// Get route parameters
	playerID := e.Vars["player_id"]
	action := e.Vars["action"]
	page, _ := strconv.Atoi(e.Vars["page"])

	e.DeferUpdateMessage()
	if playerID != e.User().ID.String() {
		return nil
	}

	switch action {
	case componentActionNext:
		page++
	case componentActionPrev:
		page--
	default:
	}

	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	username := e.User().Username
	if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	}

	embed, components := cmd.generator(username, p, page)

	_, err := e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
