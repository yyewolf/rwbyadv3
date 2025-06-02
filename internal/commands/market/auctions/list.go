package auctions

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"entgo.io/ent/dialect/sql"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/auction"
	"github.com/yyewolf/rwbyadv3/ent/auctionbid"
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
		WithBids(func(abq *ent.AuctionBidQuery) {
			abq.Order(auctionbid.ByPrice(sql.OrderDesc()))
		}).
		All(context.Background())
	if err != nil {
		return discord.Embed{}, nil, err
	}

	field.Name = fmt.Sprintf("Auctions (page %d/%d) :", page+1, maxPage+1)

	for i, auction := range auctions {
		idx := page*int(perPage) + i + 1
		field.Value += fmt.Sprintf("`N°%d | %d Ⱡ | %s`\n", idx, auction.GetPrice(), auction.Edges.Card.PartialString())
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

func (cmd *auctionsCommand) GetAuctions(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := event.User().Username
	if event.User().GlobalName != nil {
		username = *event.User().GlobalName
	}

	embed, components, err := cmd.generator(username, currentPlayer, 0)
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components),
	)
}

func (cmd *auctionsCommand) HandleGetAuctionsInteraction(logger *logrus.Entry, data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
	// Get route parameters
	playerID := event.Vars["player_id"]
	action := event.Vars["action"]
	page, _ := strconv.Atoi(event.Vars["page"])

	event.DeferUpdateMessage()
	if playerID != event.User().ID.String() {
		return nil
	}

	switch action {
	case componentActionNext:
		page++
	case componentActionPrev:
		page--
	default:
	}

	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := event.User().Username
	if event.User().GlobalName != nil {
		username = *event.User().GlobalName
	}

	embed, components, err := cmd.generator(username, currentPlayer, page)
	if err != nil {
		return utils.ComponentError(logger, event, err)
	}

	_, err = event.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddEmbeds(embed).
			AddContainerComponents(components).
			Build(),
	)
	return err
}
