package listings

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/listing"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

var (
	perPage = 10.0
)

func (cmd *listingsCommand) generator(username string, currentPlayer *ent.Player, page int) (discord.Embed, discord.ContainerComponent, error) {
	embed := discord.NewEmbedBuilder()
	embed.SetTitlef("%s's listings :", username)
	// embed.SetDescriptionf("To select a character, please use %s.", cmd.app.CommandMention("select"))
	embed.SetColor(cmd.app.Config().App.BotColor)
	embed.SetEmbedFooter(cmd.app.Footer())

	// Pagination here
	count, err := currentPlayer.QueryListings().Count(context.Background())
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

	listings, err := currentPlayer.QueryListings().
		Order(listing.ByCreateTime()).
		Offset(page * int(perPage)).
		Limit(int(perPage)).
		WithCard().
		All(context.Background())
	if err != nil {
		return discord.Embed{}, nil, err
	}

	field.Name = fmt.Sprintf("Listings (page %d/%d) :", page+1, maxPage+1)

	for i, listing := range listings {
		idx := page*int(perPage) + i + 1
		field.Value += fmt.Sprintf("`N°%d | %s`\n", idx, listing.Edges.Card.FullString())
	}

	if len(listings) == 0 {
		field.Name = "Listings :"
		field.Value = "You have no listings to be shown."
	}

	embed.AddFields(field)

	customID := fmt.Sprintf("/listings/%s/%d", currentPlayer.ID, page)

	return embed.Build(), discord.NewActionRow(
		discord.NewSecondaryButton("◀️ Prev", customID+"/"+componentActionPrev),
		discord.NewSecondaryButton("🔄 Refresh", customID+"/"+componentActionRefresh),
		discord.NewSecondaryButton("▶️ Next", customID+"/"+componentActionNext),
	), nil
}

func (cmd *listingsCommand) GetListings(logger *logrus.Entry, event *handler.CommandEvent) error {
	p := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := event.User().Username
	if event.User().GlobalName != nil {
		username = *event.User().GlobalName
	}

	embed, components, err := cmd.generator(username, p, 0)
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

func (cmd *listingsCommand) HandleGetListingsInteraction(logger *logrus.Entry, data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
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

	p := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	username := event.User().Username
	if event.User().GlobalName != nil {
		username = *event.User().GlobalName
	}

	embed, components, err := cmd.generator(username, p, page)
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
