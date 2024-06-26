package preprod

import (
	"context"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
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

			h.Command("/preprod/player_xp", builder.WithContext(
				app,
				cmd.GivePlayerXP,
				builder.WithPlayer(),
			))

			h.Command("/preprod/player_level", builder.WithContext(
				app,
				cmd.SetPlayerLevel,
				builder.WithPlayer(),
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
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "player_xp",
					Description: "Give yourself xp",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "amount",
							Description: "How much XP do you need ?",
							MinValue:    utils.Optional(0),
							Required:    true,
						},
					},
				},
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "player_level",
					Description: "Set your player level",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "level",
							Description: "What level do you want",
							MinValue:    utils.Optional(0),
							Required:    true,
						},
					},
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

func (cmd *preprodCommand) GivePlayerXP(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	amount := e.SlashCommandInteractionData().Int("amount")

	levelUp := utils.Players.GiveXP(p, int64(amount))

	p.UpdateG(context.Background(), boil.Infer())

	return e.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("Ok done, level up : %v", levelUp).
			Build(),
	)
}

func (cmd *preprodCommand) SetPlayerLevel(e *handler.CommandEvent) error {
	p := e.Ctx.Value(builder.PlayerKey).(*models.Player)

	level := e.SlashCommandInteractionData().Int("level")

	p.Level = level
	p.XP = 0
	p.NextLevelXP = utils.Players.CalculateNextLevelXP(p)

	p.UpdateG(context.Background(), boil.Infer())

	return e.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("Ok done.").
			Build(),
	)
}
