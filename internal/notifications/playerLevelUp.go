package notifications

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/ent/schema/enums"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/stats"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type PlayerLevelUpParams struct {
	Player      *ent.Player
	LevelBefore int64
}

func DispatchPlayerLevelUp(app interfaces.App, p *ent.Player, levelBefore int64) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("player_level_up_%s_%d", p.ID, p.Level),
		TaskQueue: app.Config().Temporal.TaskQueue,
	}
	app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, Repository.NotifyPlayerLevelUpWorkflow, &PlayerLevelUpParams{
		Player:      p,
		LevelBefore: levelBefore,
	})
}

type PlayerLevelUpRewards struct {
	Boxes     int
	RareBoxes int
	Liens     int64
	GoldStars int
	Backpacks int
}

func (n *NotificationsRepository) NotifyPlayerLevelUpWorkflow(ctx workflow.Context, params *PlayerLevelUpParams) error {
	// TODO : Add check for DMs, GuildChannels, and if the user wants the notification at all

	// Do rewards
	var rewards PlayerLevelUpRewards

	levelEarned := params.Player.Level - params.LevelBefore
	for i := int64(0); i < levelEarned; i++ {
		// 12.5% chance of getting lootboxes
		if stats.HasChance(12.5) {
			amount := rand.Intn(int(math.Sqrt(float64(params.Player.Level)))) + 1
			rewards.Boxes += amount
		}

		// 6.5% chance of getting rare lootboxes
		if stats.HasChance(6.5) {
			amount := rand.Intn(int(math.Sqrt(float64(params.Player.Level)))) + 1
			rewards.Boxes += amount
		}

		// Get liens every level
		rewards.Liens += rand.Int63n(153+(params.Player.Level-i+1)*6) + 54

		// Every 10 levels
		if (params.Player.Level-i+1)%10 == 0 {
			// 17.5% chance
			if stats.HasChance(17.5) {
				rewards.GoldStars++
			}

			// 10% chance
			if stats.HasChance(10) {
				rewards.Backpacks++
			}
		}
	}

	var newCtx = context.Background()

	err := ent.WithTx(newCtx, n.app.Db(), func(tx *ent.Tx) error {
		for range rewards.Boxes {
			_, err := tx.LootBox.Create().
				SetPlayerID(params.Player.ID).
				SetType(enums.LootBoxClassic).
				Save(newCtx)
			if err != nil {
				return err
			}
		}
		for range rewards.RareBoxes {
			_, err := tx.LootBox.Create().
				SetPlayerID(params.Player.ID).
				SetType(enums.LootBoxRare).
				Save(newCtx)
			if err != nil {
				return err
			}
		}

		// TODO : Add gold stars
		_, err := tx.Player.UpdateOne(params.Player).
			AddLiens(rewards.Liens).
			AddBackpackLevel(int64(rewards.Backpacks)).
			Save(newCtx)

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	c := n.app.Client()
	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(params.Player.ID))
	if err != nil {
		return err
	}

	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Congratulations !").
				SetDescriptionf(
					utils.Joinln(
						"<@%s>, you leveled up!",
						"Level : **%d**.",
						"XP : **%d/%d**",
						"",
						"You earned :",
						"%d Box(es)",
						"%d Rare Box(es)",
						"%d Ⱡ (Liens)",
						"%d Backpack(s)",
						"%d Gold Star(s)",
					),
					params.Player.ID,
					params.Player.Level,
					params.Player.ExperiencePoints, params.Player.ExperiencePointsThreshold,
					rewards.Boxes,
					rewards.RareBoxes,
					rewards.Liens,
					rewards.Backpacks,
					rewards.GoldStars,
				).
				SetColor(n.app.Config().App.BotColor).
				Build(),
		).
		Build(),
	)
	return err
}
