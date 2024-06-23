package ping

import (
	"context"
	"log"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

const (
	commandName        = "ping"
	commandDescription = "Pongs!"
)

type pingCommand struct {
	app interfaces.App
	// jobHandler interfaces.JobHandler
}

func PingCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd pingCommand

	cmd.app = app
	app.Worker().RegisterWorkflow(cmd.DelayedPongWorkflow)

	return builder.NewCommand(
		builder.WithCommandName(commandName),
		builder.WithDescription(commandDescription),
		builder.WithRegisterFunc(func(h *handler.Mux) error {
			h.Command("/"+commandName, cmd.HandleCommand)
			return nil
		}),
		builder.WithSlashCommand(discord.SlashCommandCreate{
			Name:        commandName,
			Description: commandDescription,
		}),
	)
}

func (cmd *pingCommand) HandleCommand(e *handler.CommandEvent) error {
	// Create delayed pong in 10 seconds
	workflowOptions := client.StartWorkflowOptions{
		ID:         "delayed_pong_" + e.ID().String(),
		TaskQueue:  cmd.app.Config().Temporal.TaskQueue,
		StartDelay: 5 * time.Second,
	}

	_, err := cmd.app.Temporal().ExecuteWorkflow(context.Background(), workflowOptions, cmd.DelayedPongWorkflow, e.User().ID.String())
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Pong !").
			SetEphemeral(true),
	)
}

func (cmd *pingCommand) DelayedPongWorkflow(ctx workflow.Context, userID string) (bool, error) {
	workflow.GetLogger(ctx).Info("delayed pong has been triggered", "StartTime", workflow.Now(ctx))

	c := cmd.app.Client()

	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(userID))
	if err != nil {
		return false, err
	}

	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().SetContentf("Delayed Pong ! (you used %s)", cmd.app.CommandMention("ping")).Build())
	if err != nil {
		return false, err
	}

	return true, nil
}
