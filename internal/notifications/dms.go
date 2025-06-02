package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/ent"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type SendDmParams struct {
	Player     *ent.Player
	Message    discord.MessageCreate
	Components []discord.UnmarshalComponent
}

func UnmarshalComponents(components []discord.UnmarshalComponent) []discord.ContainerComponent {
	containerComponents := make([]discord.ContainerComponent, len(components))
	for i := range components {
		containerComponents[i] = components[i].Component.(discord.ContainerComponent)
	}
	return containerComponents
}

func MarshalComponents(components []discord.ContainerComponent) []discord.UnmarshalComponent {
	unmarshalComponents := make([]discord.UnmarshalComponent, len(components))
	for i := range components {
		unmarshalComponents[i] = discord.UnmarshalComponent{
			Component: components[i],
		}
	}
	return unmarshalComponents
}

func DispatchDm(app interfaces.App, p *ent.Player, m discord.MessageCreate, components ...discord.UnmarshalComponent) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("send_dm_%s_%s", p.ID, uuid.NewString()),
		TaskQueue: app.Config().Temporal.TaskQueue,
	}

	app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, Repository.SendDmWorkflow, &SendDmParams{
		Player:     p,
		Message:    m,
		Components: components,
	})
}

func (n *NotificationsRepository) SendDmWorkflow(ctx workflow.Context, params *SendDmParams) error {
	// TODO : Add check for DMs, GuildChannels, and if the user wants the notification at all

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	var activityResult bool
	err := workflow.ExecuteActivity(ctx, n.SendDmActivity, params).Get(ctx, &activityResult)
	if err != nil {
		return err
	}

	return err
}

func (n *NotificationsRepository) SendDmActivity(ctx context.Context, params *SendDmParams) (bool, error) {
	c := n.app.Client()

	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(params.Player.ID))
	if err != nil {
		return false, err
	}

	if len(params.Components) > 0 {
		params.Message.Components = UnmarshalComponents(params.Components)
	}

	_, err = c.Rest().CreateMessage(ch.ID(), params.Message)
	if err != nil {
		logrus.WithError(err).WithField("channel_id", ch.ID()).Error("Failed to send DM message")
		return false, err
	}

	return true, nil
}
