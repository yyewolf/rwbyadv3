package builder

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type RegisterFunc func(h *handler.Mux) error

type Command struct {
	Name        string
	Description string
	menu        *Menu

	createCommand discord.SlashCommandCreate

	register RegisterFunc

	bot.Client
	*handler.Mux
}

func NewCommand(opts ...CommandOption) *Command {
	var cmd = &Command{}

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func (cmd *Command) Register() error {
	return cmd.register(cmd.Mux)
}
