package stars

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
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
	if daily.HasVoted {
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

	// Do rewards
	// money := (rand.Intn(225) + 68) * (ctx.Player.Status.DailyStreak%10 + 1)
	// earnings = append(earnings, fmt.Sprintf("**%d**Ⱡ", money))

	// var normalLootBoxes int
	// var normalGrimmBoxes int
	// // random lootbox loot
	// for i := 0; i < 3; i++ {
	// 	rng := rand.Float64() * 100
	// 	if rng < 7.5 {
	// 		t := rand.Intn(2)
	// 		if t == 0 {
	// 			ctx.Player.Boxes.Boxes++
	// 			normalLootBoxes++
	// 		} else {
	// 			ctx.Player.Boxes.GrimmBoxes++
	// 			normalGrimmBoxes++
	// 		}
	// 	}
	// }

	// if normalLootBoxes == 3 {
	// 	normalLootBoxes = 0
	// 	ctx.Player.Boxes.Boxes -= 3
	// 	ctx.Player.Boxes.RareBoxes++
	// 	earnings = append(earnings, "**1** Rare Box")
	// }
	// if normalGrimmBoxes == 3 {
	// 	normalGrimmBoxes = 0
	// 	ctx.Player.Boxes.GrimmBoxes -= 3
	// 	ctx.Player.Boxes.RareGrimmBoxes++
	// 	earnings = append(earnings, "**1** Rare Grimm Box")
	// }

	// if normalLootBoxes > 0 {
	// 	earnings = append(earnings, fmt.Sprintf("**%d** Box(es)", normalLootBoxes))
	// }
	// if normalGrimmBoxes > 0 {
	// 	earnings = append(earnings, fmt.Sprintf("**%d** Grimm Box(es)", normalGrimmBoxes))
	// }

	// cp := ctx.Player.CalcCP(0.6)
	// earnings = append(earnings, fmt.Sprintf("**%d** CP", cp))

	// content := &discordgo.MessageEmbed{
	// 	Title:       fmt.Sprintf("Daily Reward : (%d 🔥)", ctx.Player.Status.DailyStreak),
	// 	Description: "Thank you for your vote!\n\nYou earned :\n",
	// 	Thumbnail: &discordgo.MessageEmbedThumbnail{
	// 		URL: ctx.Author.AvatarURL("512"),
	// 	},
	// 	Color: config.Botcolor,
	// }

	// for _, earning := range earnings {
	// 	content.Description += fmt.Sprintf("=> %s\n", earning)
	// }

	// ctx.Reply(discord.ReplyParams{
	// 	Content: content,
	// })

	// ctx.GiveCP(cp, true)
	// ctx.Player.Boxes.Save()

	// tx, err := boil.BeginTx(e.Ctx, nil)
	// if err != nil {
	// 	return utils.CommandError(e, err)
	// }

	// daily.HasVoted = false
	// _, err = player.R.Daily.Update(
	// 	e.Ctx,
	// 	tx,
	// 	boil.Whitelist(
	// 		models.DailyColumns.HasVoted,
	// 		models.DailyColumns.LastVoteAt,
	// 		models.DailyColumns.Streak,
	// 	),
	// )
	// if err != nil {
	// 	return utils.CommandError(e, err)
	// }

	return nil
}
