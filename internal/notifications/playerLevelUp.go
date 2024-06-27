package notifications

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/stats"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type PlayerLevelUpParams struct {
	Player      *models.Player
	LevelBefore int
}

func DispatchPlayerLevelUp(app interfaces.App, p *models.Player, levelBefore int) {
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
	Liens     int
	GoldStars int
	Backpacks int
}

func (n *NotificationsRepository) NotifyPlayerLevelUpWorkflow(ctx workflow.Context, params *PlayerLevelUpParams) error {
	// TODO : Add check for DMs, GuildChannels, and if the user wants the notification at all

	// Do rewards
	var rewards PlayerLevelUpRewards

	levelEarned := params.Player.Level - params.LevelBefore
	for i := 0; i < levelEarned; i++ {
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
		rewards.Liens += rand.Intn(153+(params.Player.Level-i+1)*6) + 54

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

	tx, err := boil.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	for range rewards.Boxes {
		params.Player.AddLootBoxes(context.Background(), tx, true, &models.LootBox{
			ID:       uuid.NewString(),
			PlayerID: params.Player.ID,
			Type:     models.LootBoxesTypeClassic,
		})
	}
	for range rewards.RareBoxes {
		params.Player.AddLootBoxes(context.Background(), tx, true, &models.LootBox{
			ID:       uuid.NewString(),
			PlayerID: params.Player.ID,
			Type:     models.LootBoxesTypeRare,
		})
	}

	err = params.Player.Reload(context.Background(), tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	params.Player.Liens += int64(rewards.Liens)
	params.Player.BackpackLevel += rewards.Backpacks
	// TODO : Add gold stars

	_, err = params.Player.Update(context.Background(), tx, boil.Whitelist(
		models.PlayerColumns.Liens,
		models.PlayerColumns.BackpackLevel,
	))
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
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
					params.Player.XP, params.Player.NextLevelXP,
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
