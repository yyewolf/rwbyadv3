package ping

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

const (
	commandName        = "ping"
	commandDescription = "Pongs!"
)

type pingCommand struct {
	app        interfaces.App
	jobHandler interfaces.JobHandler
}

func PingCommand(ms *builder.MenuStore, app interfaces.App) *builder.Command {
	var cmd pingCommand

	cmd.app = app
	cmd.jobHandler = app.JobHandler()
	cmd.jobHandler.RegisterJobKey("delayed_pong", cmd.DelayedPongJob)

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
	_, err := cmd.jobHandler.ScheduleJob(
		"delayed_pong",
		e.ID().String(),
		time.Now().Add(10*time.Second),
		map[string]interface{}{"user_id": e.User().ID.String()},
	)
	if err != nil {
		logrus.WithError(err).Error("failed to schedule delayed pong job")
	}

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetContentf("Pong !").
			SetEphemeral(true),
	)
}

func (cmd *pingCommand) DelayedPongJob(params map[string]interface{}) error {
	// get user id from param "user_id"
	id := params["user_id"].(string)

	c := cmd.app.Client()

	ch, err := c.Rest().CreateDMChannel(snowflake.MustParse(id))
	if err != nil {
		return err
	}

	_, err = c.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().SetContentf("Delayed Pong ! (you used %s)", cmd.app.CommandMention("ping")).Build())

	return err
}
