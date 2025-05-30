package trades

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/card"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/notifications"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/web/api"
)

type CreateTradeRequest struct {
	Offer   []string `json:"offer"`
	Receive []string `json:"receive"`
}

func stringListToUUIDList(strs []string) ([]uuid.UUID, error) {
	var uuids []uuid.UUID
	for _, str := range strs {
		uuid, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, uuid)
	}
	return uuids, nil
}

func discordTradeMessage(offer []*ent.Card, receive []*ent.Card, recipient bool) (discord.MessageCreate, []discord.UnmarshalComponent) {
	var fields []discord.EmbedField
	var components []discord.UnmarshalComponent

	if len(offer) > 0 {
		offerField := discord.EmbedField{
			Name: "You will lose :",
		}
		for _, card := range offer {
			offerField.Value += fmt.Sprintf("`%s`", card.FullString()) + "\n"
		}
		fields = append(fields, offerField)
	}

	if len(receive) > 0 {
		receiveField := discord.EmbedField{
			Name: "You will receive :",
		}
		for _, card := range receive {
			receiveField.Value += fmt.Sprintf("`%s`", card.FullString()) + "\n"
		}
		fields = append(fields, receiveField)
	}

	title := "Trade Request Sent"
	description := "You have sent a trade request"
	if recipient {
		title = "Trade Request Received"
		description = "You have received a trade request"
	}

	msg := discord.NewMessageCreateBuilder().
		AddEmbeds(
			discord.NewEmbedBuilder().
				SetTitle(title).
				SetDescription(description).
				SetColor(0x00FF00).
				AddFields(fields...).
				Build(),
		)

	if recipient {
		row := discord.NewActionRow(
			discord.NewPrimaryButton("Accept Trade", "accept_trade").
				WithEmoji(discord.ComponentEmoji{Name: "✅"}),
			discord.NewDangerButton("Decline Trade", "decline_trade").
				WithEmoji(discord.ComponentEmoji{Name: "❎"}),
		)

		components = notifications.MarshalComponents([]discord.ContainerComponent{row})
	}

	return msg.Build(), components
}

func (h *TradeApiHandler) CreateTrade(app interfaces.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := utils.GetSessionFromContext(c)
		initiator := session.Edges.Player

		recipient, err := app.Db().Player.Query().
			Where(player.ID(c.Param("playerId"))).
			Only(c.Request().Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return HandleErrorJson(c, err, "Player not found")
			}
			return HandleErrorJson(c, err, "Failed to fetch player")
		}

		initiatorTrades, err := initiator.QueryTrades().All(c.Request().Context())
		if err != nil {
			return HandleErrorJson(c, err, "Failed to fetch initiator's trades")
		}

		if len(initiatorTrades) >= 5 {
			return HandleErrorJson(c, nil, "You have reached the maximum number of trades (5)")
		}

		for _, trade := range initiatorTrades {
			if trade.ReceiverID == recipient.ID {
				return HandleErrorJson(c, nil, "You already have a trade with this player")
			}
		}

		var req CreateTradeRequest
		if err := c.Bind(&req); err != nil {
			return HandleErrorJson(c, err, "Failed to bind request")
		}

		initatorCardsUuids, err := stringListToUUIDList(req.Offer)
		if err != nil {
			return HandleErrorJson(c, err, "Invalid card UUID format in offer")
		}
		receiveCardsUuids, err := stringListToUUIDList(req.Receive)
		if err != nil {
			return HandleErrorJson(c, err, "Invalid card UUID format in receive")
		}

		// Check initator's cards
		initatorCards, err := initiator.QueryCards().
			Where(card.IDIn(initatorCardsUuids...)).
			All(c.Request().Context())
		if err != nil {
			return HandleErrorJson(c, err, "Failed to fetch initiator's cards")
		}

		if len(initatorCards) != len(initatorCardsUuids) {
			return HandleErrorJson(c, nil, "Some cards in the offer do not belong to you")
		}

		// Check recipient's cards
		recipientCards, err := recipient.QueryCards().
			Where(card.IDIn(receiveCardsUuids...)).
			All(c.Request().Context())
		if err != nil {
			return HandleErrorJson(c, err, "Failed to fetch recipient's cards")
		}

		if len(recipientCards) != len(receiveCardsUuids) {
			return HandleErrorJson(c, nil, "Some cards in the receive do not belong to the recipient")
		}

		_, err = h.app.Db().Trade.Create().
			SetInitiator(initiator).
			SetReceiver(recipient).
			SetOfferCards(req.Offer).
			SetReceiveCards(req.Receive).
			Save(c.Request().Context())
		if err != nil {
			return HandleErrorJson(c, err, "Failed to create trade")
		}

		msg, components := discordTradeMessage(initatorCards, recipientCards, false)
		recipientMsg, recipientComponents := discordTradeMessage(recipientCards, initatorCards, true)

		// Send Discord messages to both players
		notifications.DispatchDm(h.app, initiator, msg, components...)
		notifications.DispatchDm(h.app, recipient, recipientMsg, recipientComponents...)

		return api.OkRedirect(c, "/landing/discord/")
	}
}
