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
			h.Command("/open", builder.WithContext(
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerLootBoxes(),
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

	if len(p.R.LootBoxes) == 0 {
		return e.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("You don't have any lootboxes :(").
			SetEphemeral(true).
			Build(),
		)
	}

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	lootBox := p.R.LootBoxes[0]

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

	tx.Commit()

	return e.CreateMessage(discord.NewMessageCreateBuilder().
		SetContentf("Looted `%s`, you have %d loot boxes left.", utils.Cards.FullString(c), len(p.R.LootBoxes)-1).
		Build(),
	)
}
