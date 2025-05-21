package daily

import (
	"math/rand"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/pkg/loottables"
	"github.com/yyewolf/rwbyadv3/pkg/loottables/item"
)

const (
	commandName        = "daily"
	commandDescription = "Get a daily reward for voting on Top.GG"
)

type dailyCommand struct {
	app interfaces.App
}

func DailyCommand(menus *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd dailyCommand

	cmd.app = app

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, builder.WithContext(
				app,
				cmd.HandleCommand,
				builder.WithPlayer(),
				builder.WithPlayerDaily(),
			))
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *dailyCommand) HandleCommand(logger *logrus.Entry, event *handler.CommandEvent) error {
	currentPlayer := event.Ctx.Value(builder.NewPlayerKey).(*ent.Player)
	daily := currentPlayer.Edges.Daily
	if !daily.HasVoted {
		return event.Respond(
			discord.InteractionResponseTypeCreateMessage,
			discord.NewMessageCreateBuilder().
				SetEmbeds(
					discord.NewEmbedBuilder().
						SetTitle("Daily Reward :").
						SetColor(cmd.app.Config().App.BotColor).
						SetDescriptionf("You did not vote yet, you can click [here](%s) to vote !", cmd.app.Config().TopGg.Url).
						SetEmbedFooter(cmd.app.Footer()).
						Build(),
				).
				SetEphemeral(true).
				Build(),
		)
	}

	var random = rand.New(rand.NewSource(rand.Int63()))

	// Create the loot table
	var lootTable = loottables.New(
		item.New(&loots.Liens{}, 10).
			OnlyDropOnce().
			WithAmountRange(100, 1500, 1).
			WithRepartitionFunc(
				item.RepartitionGaussian(int(120*(daily.Streak%7+1)), 50),
			),
		item.New(item.Nothing{}, 20),
	)

	list := lootTable.ChooseRandomItems(random, 5)

	texts := make([]string, 0)

	err := ent.WithTx(event.Ctx, cmd.app.Db(), func(tx *ent.Tx) error {
		// Give loots and create text
		for _, loot := range list {
			if loot, ok := loot.(loots.Loot); ok {
				loot.PickedUp(tx, currentPlayer)
				texts = append(texts, loot.RewardText(list))
			}
		}

		err := tx.Daily.UpdateOne(daily).
			SetHasVoted(false).
			Exec(event.Ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return utils.CommandError(logger, event, err)
	}

	return event.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEmbeds(
				discord.NewEmbedBuilder().
					SetTitlef("Daily Reward : (%d 🔥)", daily.Streak).
					SetColor(cmd.app.Config().App.BotColor).
					SetDescriptionf(
						utils.Joinln(
							"Thank you for your vote!",
							"Here's what you got :",
							"",
							utils.Joinln(texts...),
						),
					).
					SetEmbedFooter(cmd.app.Footer()).
					Build(),
			).
			SetEphemeral(true).
			Build(),
	)
}
