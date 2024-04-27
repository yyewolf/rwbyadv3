package bugs

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type BugsMenu struct {
	commands []interfaces.Command

	cr interfaces.CommandRepository
}

func New(cr interfaces.CommandRepository) interfaces.Menu {
	m := &BugsMenu{
		cr: cr,
	}

	m.commands = append(m.commands, newReportCommand(m))

	return m
}

func (m *BugsMenu) GetName() string {
	return "Bugs"
}

func (m *BugsMenu) GetEmoji() discord.Emoji {
	return discord.Emoji{
		Name: "🐞",
	}
}

func (m *BugsMenu) GetSubcommands() []interfaces.Command {
	return m.commands
}

func (m *BugsMenu) RegisterCommands() ([]api.CreateCommandData, error) {
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
