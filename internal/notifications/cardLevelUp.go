package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/internal/temporal"
	"github.com/yyewolf/rwbyadv3/internal/utils"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type CardLevelUpParams struct {
	Player *ent.Player
	Card   *ent.Card
}

func DispatchCardLevelUp(app interfaces.App, p *ent.Player, c *ent.Card) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("card_level_up_%s_%d", c.ID, c.Level),
		TaskQueue: app.Config().Temporal.TaskQueue,
	}
	app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, Repository.NotifyCardLevelUpWorkflow, &CardLevelUpParams{
		Player: p,
		Card:   c,
	})
}

func (n *NotificationsRepository) NotifyCardLevelUpWorkflow(ctx workflow.Context, params *CardLevelUpParams) error {

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var activityResult temporal.AuctionEndStatus
	err := workflow.ExecuteActivity(ctx, n.NotifyCardLevelUpActivity, params).Get(ctx, &activityResult)
	if err != nil {
		return err
	}

	return nil
}

func (n *NotificationsRepository) NotifyCardLevelUpActivity(ctx context.Context, params *CardLevelUpParams) error {
	// TODO : Add check for DMs, GuildChannels, and if the user wants the notification at all

	c := n.app.Client()
	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(params.Player.ID))
	if err != nil {
		return err
	}
	primitive := params.Card.Primitive()
	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().
		SetEmbeds(
			discord.NewEmbedBuilder().
				SetTitle("Congratulations !").
				SetDescriptionf(
					utils.Joinln(
						"<@%s>, your **%s** has leveled up!",
						"Level : **%d**.",
						"XP : **%d/%d**",
					),
					params.Player.ID, primitive.Name,
					params.Card.Level,
					params.Card.ExperiencePoints, params.Card.ExperiencePointsThreshold,
				).
				SetThumbnail(utils.Cards.IconURI(params.Card)).
				SetColor(n.app.Config().App.BotColor).
				Build(),
		).
		Build(),
	)
	return err
}
