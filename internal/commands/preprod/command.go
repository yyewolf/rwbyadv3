package preprod

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/utils"
)

const (
	commandName        = "preprod"
	commandDescription = "All commands available to the preprod environmet"
)

type preprodCommand struct {
	app interfaces.App
}

func PreprodCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
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

			h.Command("/preprod/liens", builder.WithContext(
				app,
				cmd.SetLiens,
				builder.WithPlayer(),
			))

			h.Command("/preprod/backpacks", builder.WithContext(
				app,
				cmd.SetLiens,
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
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "liens",
					Description: "Set your amount of liens",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "liens",
							Description: "What amount of liens do you want",
							MinValue:    utils.Optional(0),
							Required:    true,
						},
					},
				},
				&discord.ApplicationCommandOptionSubCommand{
					Name:        "backpacks",
					Description: "Set your amount of backpacks",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "backpacks",
							Description: "What amount of backpacks do you want",
							MinValue:    utils.Optional(1),
							Required:    true,
						},
					},
				},
			},
		}),
	)
}

func (cmd *preprodCommand) DoLootBox(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	for _, t := range enums.LootBoxTypes() {
		cmd.app.Db().LootBox.Create().
			SetPlayerID(currentPlayer.ID).
			SetType(t).
			Save(event.Ctx)
	}

	count, err := currentPlayer.QueryLootboxes().Count(event.Ctx)
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("You now have %d loot boxes", count).
			Build(),
	)
}

func (cmd *preprodCommand) GivePlayerXP(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)
	amount := event.SlashCommandInteractionData().Int("amount")
	levelUp := currentPlayer.GiveXP(int64(amount))

	currentPlayer.Update().
		SetExperiencePoints(currentPlayer.ExperiencePoints).
		SetExperiencePointsThreshold(currentPlayer.ExperiencePointsThreshold).
		SetLevel(currentPlayer.Level).
		Save(event.Ctx)

	return event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("Ok done, level up : %v", levelUp).
			Build(),
	)
}

func (cmd *preprodCommand) SetPlayerLevel(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)
	level := event.SlashCommandInteractionData().Int("level")
	currentPlayer.Level = int64(level)
	currentPlayer.ExperiencePoints = 0
	currentPlayer.ExperiencePointsThreshold = currentPlayer.GetNextLevelXP()

	currentPlayer.Update().
		SetExperiencePoints(currentPlayer.ExperiencePoints).
		SetExperiencePointsThreshold(currentPlayer.ExperiencePointsThreshold).
		SetLevel(currentPlayer.Level).
		Save(event.Ctx)

	return event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("Ok done.").
			Build(),
	)
}

func (cmd *preprodCommand) SetLiens(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)

	currentPlayer.Update().
		SetLiens(int64(event.SlashCommandInteractionData().Int("liens"))).
		Save(event.Ctx)

	return event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("Ok done.").
			Build(),
	)
}

func (cmd *preprodCommand) SetBackpack(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)
	backpacks := event.SlashCommandInteractionData().Int("backpacks")
	currentPlayer.Update().SetBackpackLevel(int64(backpacks)).Save(event.Ctx)

	return event.CreateMessage(
		discord.NewMessageCreateBuilder().
			SetContentf("Ok done.").
			Build(),
	)
}
