package system

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type SystemMenu struct {
	commands []interfaces.Command

	cr interfaces.CommandRepository
}

func New(cr interfaces.CommandRepository) interfaces.Menu {
	m := &SystemMenu{
		cr: cr,
	}

	m.commands = append(m.commands, newPingCommand(m))

	return m
}

func (m *SystemMenu) GetName() string {
	return "System"
}

func (m *SystemMenu) GetEmoji() discord.Emoji {
	return discord.Emoji{
		Name: "🔧",
	}
}

func (m *SystemMenu) GetSubcommands() []interfaces.Command {
	return m.commands
}

func (m *SystemMenu) RegisterCommands() ([]api.CreateCommandData, error) {
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
