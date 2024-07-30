package daily

import (
	"context"
	"math/rand"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
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

func DailyCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
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

func (cmd *dailyCommand) HandleCommand(e *handler.CommandEvent) error {
	player := e.Ctx.Value(builder.PlayerKey).(*models.Player)
	daily := player.R.GetDaily()
	if !daily.HasVoted {
		return e.Respond(
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
		item.New(&loots.Liens{}, 10).OnlyDropOnce().WithAmountRange(100, 1500, 1).WithRepartitionFunc(item.RepartitionGaussian(120*(daily.Streak%7+1), 50)),
		item.New(item.Nothing{}, 20),
	)

	list := lootTable.ChooseRandomItems(random, 5)

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return utils.CommandError(e, err)
	}

	// Give loots and create text
	texts := make([]string, 0)
	for _, loot := range list {
		if loot, ok := loot.(loots.Loot); ok {
			loot.PickedUp(tx, player)
			texts = append(texts, loot.RewardText(list))
		}
	}

	// Remove capability from player to claim again
	daily.HasVoted = false
	_, err = player.R.Daily.Update(
		e.Ctx,
		tx,
		boil.Whitelist(
			models.DailyColumns.HasVoted,
		),
	)
	if err != nil {
		tx.Rollback()
		return utils.CommandError(e, err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.CommandError(e, err)
	}

	return e.Respond(
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
