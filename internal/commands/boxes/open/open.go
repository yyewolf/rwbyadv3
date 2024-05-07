package open

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/cards"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/lootbox"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "open"
	commandDescription = "Open a lootbox (still in demo)"

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
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerLootBoxes(),
			))

			h.ButtonComponent(componentId, builder.WithContextD(
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
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	components := cmd.generator(p)

	return e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContainerComponents(components).
		Build(),
	)
}

func (cmd *openCommand) HandleInteraction(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
	// Get route parameters
	playerID := e.Vars["player_id"]
	boxType := models.LootBoxesType(e.Vars["box_type"])

	e.DeferUpdateMessage()
	if playerID != e.User().ID.String() {
		return nil
	}

	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	lootBox, found := utils.Players.TakeFirstLootBoxOf(p, boxType)
	if !found {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("You don't have any more of these loot boxes :(").
			SetEphemeral(true).
			Build(),
		)
	}

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = lootBox.Delete(context.Background(), tx, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	// var c *models.Card

	// switch lootBox.Type {
	// case models.LootBoxesTypeClassic:
	// case models.LootBoxesTypeRare:
	// case models.LootBoxesTypeLimited:
	// case models.LootBoxesTypeSpecial:
	// }

	c := lootbox.NormalLootBox.PickCard(cards.Cards)
	c.PlayerID = e.User().ID.String()

	err = c.Insert(context.Background(), tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return err
	}

	err = c.SetCardsStat(context.Background(), tx, true, utils.Cards.GenerateStats(c))
	if err != nil {
		tx.Rollback()
		return err
	}

	// If user does not have a selected card, this is it :
	if p.R.SelectedCard == nil {
		err = p.SetSelectedCard(context.Background(), tx, false, c)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	f, embed, _ := utils.Cards.Message(c)

	_, err = e.CreateFollowupMessage(discord.NewMessageCreateBuilder().
		SetFiles(f).
		SetEmbeds(embed).
		Build(),
	)

	utils.Players.DeleteBoxFromPlayer(p, lootBox)
	components := cmd.generator(p)
	e.UpdateInteractionResponse(
		discord.NewMessageUpdateBuilder().
			AddContainerComponents(components).
			Build(),
	)
	return err
}
