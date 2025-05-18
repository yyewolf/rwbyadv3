package utils

import (
	"runtime/debug"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/builder"
)

func CommandError(e *handler.CommandEvent, err error) error {
	logrus.
		WithField(string(builder.ContextIdKey), e.Ctx.Value(builder.ContextIdKey)).
		WithError(err).
		WithField("stack", string(debug.Stack())).
		Error("An error occurred while handling a command")

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEphemeral(true).
			SetContentf("Uh-oh! An error occurred, if this persists, please contact support with the following code: %s", e.Ctx.Value(builder.ContextIdKey)),
	)
}

func ComponentError(e *handler.ComponentEvent, err error) error {
	logrus.
		WithField(string(builder.ContextIdKey), e.Ctx.Value(builder.ContextIdKey)).
		WithError(err).
		WithField("stack", string(debug.Stack())).
		Error("An error occurred while handling a component")

	return e.Respond(
		discord.InteractionResponseTypeCreateMessage,
		discord.NewMessageCreateBuilder().
			SetEphemeral(true).
			SetContentf("Uh-oh! An error occurred, if this persists, please contact support with the following code: %s", e.Ctx.Value(builder.ContextIdKey)),
	)
}
