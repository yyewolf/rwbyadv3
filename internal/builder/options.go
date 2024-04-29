package builder

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type CommandOption func(a *Command)
type MenuOption func(a *Menu)

func WithCommandName(name string) CommandOption {
	return func(cmd *Command) {
		cmd.Name = name
	}
}

func WithDescription(desc string) CommandOption {
	return func(cmd *Command) {
		cmd.Description = desc
	}
}

func WithApp(app interfaces.App) CommandOption {
	return func(cmd *Command) {
		cmd.Client = app.Client()
		cmd.Mux = app.Handler()
	}
}

func WithRegisterFunc(r RegisterFunc) CommandOption {
	return func(cmd *Command) {
		cmd.register = r
	}
}
func WithSlashCommand(r discord.SlashCommandCreate) CommandOption {
	return func(cmd *Command) {
		cmd.createCommand = r
	}
}

func WithMenuName(name string) MenuOption {
	return func(menu *Menu) {
		menu.Name = name
	}
}

func WithEmoji(emoji discord.Emoji) MenuOption {
	return func(menu *Menu) {
		menu.Emoji = emoji
	}
}

func WithCommands(cmds ...*Command) MenuOption {
	return func(menu *Menu) {
		menu.Commands = append(menu.Commands, cmds...)

		for _, cmd := range cmds {
			cmd.Mux = menu.h
			cmd.menu = menu
		}
	}
}
