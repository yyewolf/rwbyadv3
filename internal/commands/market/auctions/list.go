package auctions

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

var (
	perPage = 10.0
)

func (cmd *auctionsCommand) generator(username string, p *ent.Player, page int) (discord.Embed, discord.ContainerComponent, error) {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's auctions :", username)
	// embed.SetDescriptionf("To select a character, please use %s.", cmd.app.CommandMention("select"))
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	count, err := p.QueryAuctions().Count(context.Background())
	if err != nil {
		return discord.Embed{}, nil, err
	}

	maxPage := int(math.Ceil(float64(count)/perPage)) - 1

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

	auctions, err := p.QueryAuctions().
		Order(auction.ByEndsAt()).
		Offset(page * int(perPage)).
		Limit(int(perPage)).
		WithCard().
		All(context.Background())
	if err != nil {
		return discord.Embed{}, nil, err
	}

	field.Name = fmt.Sprintf("Auctions (page %d/%d) :", page+1, maxPage+1)

	for i, auction := range auctions {
		idx := page*int(perPage) + i + 1
		field.Value += fmt.Sprintf("`N°%d | %s`\n", idx, auction.Edges.Card.FullString())
	}

	if len(auctions) == 0 {
		field.Name = "Auctions :"
		field.Value = "You have no auctions to be shown."
	}

	embed.AddFields(field)

	return embed.Build(), discord.NewActionRow(
		discord.NewSecondaryButton("◀️ Prev", fmt.Sprintf(componentFormat, p.ID, page, componentActionPrev)),
		discord.NewSecondaryButton("🔄 Refresh", fmt.Sprintf(componentFormat, p.ID, page, componentActionRefresh)),
		discord.NewSecondaryButton("▶️ Next", fmt.Sprintf(componentFormat, p.ID, page, componentActionNext)),
	), nil
}

func (cmd *auctionsCommand) GetAuctions(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := e.User().Username
	if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	}

	embed, components, err := cmd.generator(username, p, 0)
	if err != nil {
		return utils.CommandError(e, err)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *auctionsCommand) HandleGetAuctionsInteraction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
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

	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := e.User().Username
	if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	}

	embed, components, err := cmd.generator(username, p, page)
	if err != nil {
		return utils.ComponentError(e, err)
	}

	_, err = e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
