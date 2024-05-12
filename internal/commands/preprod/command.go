package preprod

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
)

const (
	commandName        = "preprod"
	commandDescription = "All commands available to the preprod environmet"
)

type preprodCommand struct {
	app interfaces.App
}

func PreprodCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd preprodCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {

			h.Command("/preprod/lootbox", builder.WithContext(
				app,
				cmd.DoLootBox,
				builder.WithPlayer(),
				builder.WithPlayerLootBoxes(),
			))

			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
			// subcommands for bugs and issues
			Options: []discord.ApplicationCommandOption{
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "lootbox",
					Description: "Report a bug",
				},
			},
		}),
	)
}

func (cmd *preprodCommand) DoLootBox(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	for _, t := range models.AllLootBoxesType() {
		p.AddLootBoxesG(context.Background(), true, &models.LootBox{
			ID:       uuid.NewString(),
			PlayerID: p.ID,
			Type:     t,
		})
	}

	return e.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("You now have %d loot boxes", len(p.R.LootBoxes)).
			Build(),
	)
}
