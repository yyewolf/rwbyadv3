package open

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/ent"
	entLootbox "github.com/yyewolf/rwbyadv3/ent/lootbox"
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

func OpenCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
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

func (cmd *openCommand) HandleCommand(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	components := cmd.generator(p)

	return e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContainerComponents(components).
		Build(),
	)
}

func (cmd *openCommand) HandleInteraction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	// Get route parameters
	playerID := e.Vars["player_id"]
	boxType := enums.LootBoxType(e.Vars["box_type"])

	if playerID != e.User().ID.String() {
		return nil
	}

	p := e.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	currentLootbox, err := p.QueryLootboxes().
		Where(
			entLootbox.TypeEQ(boxType),
		).
		First(e.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return e.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("You don't have any more of these loot boxes :(").
				SetEphemeral(true).
				Build(),
			)
		}
		return utils.ComponentError(e, err)
	}

	// Check for available slots
	if utils.Players.NewAvailableSlots(p) == 0 {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("You don't have any available slots in your backpack :(").
			SetEphemeral(true).
			Build(),
		)
	}

	e.DeferUpdateMessage()

	var c *ent.Card

	err = utils.WithTx(e.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		err = tx.LootBox.DeleteOneID(currentLootbox.ID).Exec(e.Ctx)
		if err != nil {
			return err
		}

		switch currentLootbox.Type {
		case enums.LootBoxClassic:
			c = lootbox.NormalLootBox.PickCard(cards.Cards)
		case enums.LootBoxRare:
			c = lootbox.RareLootBox.PickCard(cards.Cards)
		case enums.LootBoxLimited:
			c = lootbox.LimitedLootBox.PickCard(cards.Cards)
		case enums.LootBoxSpecial:
			c = lootbox.SpecialLootBox.PickCard(cards.Cards)
		}

		c.PlayerID = e.User().ID.String()

		count, err := p.QueryCards().Count(e.Ctx)
		if err != nil {
			return err
		}

		c, err = tx.Card.Create().
			SetLevel(c.Level).
			SetPlayerID(p.ID).
			SetCardType(c.CardType).
			SetIndividualValue(c.IndividualValue).
			SetRarity(c.Rarity).
			SetPosition(float64(count)).
			SetExperiencePointsThreshold(c.ExperiencePointsThreshold).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		stats := c.GenerateStats()

		cardStats, err := tx.CardStats.Create().
			SetCardID(c.ID).
			SetArmor(stats.Armor).
			SetDamage(stats.Damage).
			SetHealing(stats.Healing).
			SetHealth(stats.Health).
			SetSpeed(stats.Speed).
			Save(e.Ctx)
		if err != nil {
			return err
		}

		c.Edges.Stats = cardStats

		// If user does not have a selected card, this is it :
		if p.Edges.SelectedCard == nil {
			_, err = tx.Player.UpdateOne(p).
				SetSelectedCardID(c.ID).
				Save(e.Ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return utils.ComponentError(e, err)
	}

	f, embed, _ := c.Message()
	embed.Footer = cmd.app.Footer()

	_, err = e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		SetFiles(f).
		SetEmbeds(embed).
		Build(),
	)

	components := cmd.generator(p)
	e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddContainerComponents(components).
			Build(),
	)
	return err
}
