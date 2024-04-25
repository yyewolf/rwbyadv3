package general

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type GeneralMenu struct {
	commands []interfaces.Command

	cr interfaces.CommandRepository
}

func New(cr interfaces.CommandRepository) interfaces.Menu {
	m := &GeneralMenu{
		cr: cr,
	}

	m.commands = append(m.commands, newHelpCommand(m))

	return m
}

func (m *GeneralMenu) GetName() string {
	return "General"
}

func (m *GeneralMenu) GetEmoji() discord.Emoji {
	return discord.Emoji{
		Name: "💻",
	}
}

func (m *GeneralMenu) GetSubcommands() []interfaces.Command {
	return m.commands
}

func (m *GeneralMenu) RegisterCommands() ([]api.CreateCommandData, error) {
	var data []api.CreateCommandData

	for _, c := range m.commands {
		val, err := c.RegisterCommand()
		if err != nil {
			return nil, err
		}
		if val == nil {
			continue
		}
		data = append(data, *val)
	}

	return data, nil
}
