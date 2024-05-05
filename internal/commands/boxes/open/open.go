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
			h.Command("/open", cmd.HandleCommand())
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *openCommand) HandleCommand() handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		c := lootbox.NormalLootBox.PickCard(cards.Cards)
		t := utils.Cards.Template(c)

		c.PlayerID = e.User().ID.String()

		tx, err := boil.BeginTx(context.Background(), nil)
		if err != nil {
			return err
		}

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
			SetContentf("Looted %s (%v, %v)", t.Name, c, c.R.CardsStat).
			Build(),
		)
	}
}
