package notifications

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type SendDmParams struct {
	Player  *models.Player
	Message discord.MessageCreate
}

func DispatchDm(app interfaces.App, p *models.Player, m discord.MessageCreate) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("send_dm_%s_%s", p.ID, uuid.NewString()),
		TaskQueue: app.Config().Temporal.TaskQueue,
	}
	app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, Repository.SendDmWorkflow, &SendDmParams{
		Player:  p,
		Message: m,
	})
}

func (n *NotificationsRepository) SendDmWorkflow(ctx workflow.Context, params *SendDmParams) error {
	// TODO : Add check for DMs, GuildChannels, and if the user wants the notification at all

	c := n.app.Client()

	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(params.Player.ID))
	if err != nil {
		return err
	}

	_, err = c.Rest().CreateMessage(ch.ID(), params.Message)
	return err
}
