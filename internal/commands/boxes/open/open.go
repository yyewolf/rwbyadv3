package open

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	entLootbox "github.com/yyewolf/rwbyadv3/ent/lootbox"
	"github.com/yyewolf/rwbyadv3/ent/player"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/lootbox"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

const (
	commandName        = "open"
	commandDescription = "Open a lootbox"

	componentId    = "/open/{player_id}/{box_type}"
	componentIdFmt = "/open/%s/%s"
)

type openCommand struct {
	app interfaces.App
}

func OpenCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd openCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerLootBoxes(),
			))

			h.ButtonComponent(componentId, builder.WithContextD(
				app,
				cmd.HandleInteraction,
				builder.WithPlayer(),
				builder.WithPlayerLootBoxes(),
				builder.WithPlayerSelectedCard(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *openCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	components := cmd.generator(currentPlayer)

	return event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContainerComponents(components).
		Build(),
	)
}

func (cmd *openCommand) HandleInteraction(logger *logrus.Entry, data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
	// Get route parameters
	playerID := event.Vars["player_id"]
	boxType := enums.LootBoxType(event.Vars["box_type"])

	if playerID != event.User().ID.String() {
		return nil
	}

	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	currentLootbox, err := currentPlayer.QueryLootboxes().
		Where(
			entLootbox.TypeEQ(boxType),
		).
		First(event.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("You don't have any more of these loot boxes :(").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.ComponentError(logger, event, err)
	}

	// Check for available slots
	if utils.Players.AvailableSlots(currentPlayer) == 0 {
		return event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("You don't have any available slots in your backpack :(").
			SetEphemeral(true).
			Build(),
		)
	}

	event.DeferUpdateMessage()

	var newCard *ent.Card

	err = ent.WithTx(event.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		err = tx.LootBox.DeleteOneID(currentLootbox.ID).Exec(event.Ctx)
		if err != nil {
			return err
		}

		switch currentLootbox.Type {
		case enums.LootBoxClassic:
			newCard = lootbox.NormalLootBox.PickCard(cards.Cards)
		case enums.LootBoxRare:
			newCard = lootbox.RareLootBox.PickCard(cards.Cards)
		case enums.LootBoxLimited:
			newCard = lootbox.LimitedLootBox.PickCard(cards.Cards)
		case enums.LootBoxSpecial:
			newCard = lootbox.SpecialLootBox.PickCard(cards.Cards)
		}

		newCard.PlayerID = event.User().ID.String()

		count, err := currentPlayer.QueryCards().Count(event.Ctx)
		if err != nil {
			return err
		}

		newCard, err = tx.Card.Create().
			SetLevel(newCard.Level).
			SetPlayerID(currentPlayer.ID).
			SetCardType(newCard.CardType).
			SetIndividualValue(newCard.IndividualValue).
			SetRarity(newCard.Rarity).
			SetPosition(float64(count)).
			SetExperiencePointsThreshold(newCard.ExperiencePointsThreshold).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		stats := newCard.GenerateStats()

		cardStats, err := tx.CardStats.Create().
			SetCardID(newCard.ID).
			SetArmor(stats.Armor).
			SetDamage(stats.Damage).
			SetHealing(stats.Healing).
			SetHealth(stats.Health).
			SetSpeed(stats.Speed).
			Save(event.Ctx)
		if err != nil {
			return err
		}

		newCard.Edges.Stats = cardStats

		// If user does not have a selected card, this is it :
		if currentPlayer.Edges.SelectedCard == nil {
			_, err = tx.Player.UpdateOne(currentPlayer).
				SetSelectedCardID(newCard.ID).
				Save(event.Ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return utils.ComponentError(logger, event, err)
	}

	logrus.WithField("card_string", newCard.FullString()).Info("new card created from lootbox")

	embedFile, embed, _ := newCard.Message()
	embed.Footer = cmd.app.Footer()

	_, err = event.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		SetFiles(embedFile).
		SetEmbeds(embed).
		Build(),
	)
	if err != nil {
		logrus.WithError(err).Error("error sending followup message")
	}

	currentPlayer, err = cmd.app.Db().Player.Query().
		Where(player.ID(currentPlayer.ID)).
		WithLootboxes().
		Only(event.Ctx)
	if err == nil {
		event.UpdateInteractionResponse(
			discord.NewMessageUpdateBuilder().
				AddContainerComponents(cmd.generator(currentPlayer)).
				Build(),
		)
	}

	return err
}
