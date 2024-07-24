package builder

import (
	"fmt"

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

type SubCommandOptionWithPrefix struct {
	Prefix string
	*discord.ApplicationCommandOptionSubCommand
}

func getSubCommandsRecursive(prefix string, opt discord.ApplicationCommandOption, out *[]SubCommandOptionWithPrefix) {
	switch o := opt.(type) {
	case *discord.ApplicationCommandOptionSubCommandGroup:
		for _, nopt := range o.Options {
			getSubCommandsRecursive(fmt.Sprintf("%s %s", prefix, o.Name), nopt, out)
		}
	case *discord.ApplicationCommandOptionSubCommand:
		*out = append(*out, SubCommandOptionWithPrefix{
			Prefix:                             prefix,
			ApplicationCommandOptionSubCommand: o,
		})
	}
}

func (cmd *Command) GetSubCommands() (out []SubCommandOptionWithPrefix) {
	for _, opt := range cmd.createCommand.Options {
		getSubCommandsRecursive(cmd.Name, opt, &out)
	}

	return out
}
