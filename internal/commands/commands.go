package commands

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/commands/bugs"
	"github.com/yyewolf/rwbyadv3/internal/commands/general"
	"github.com/yyewolf/rwbyadv3/internal/commands/system"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
)

type CommandRepositoryImpl struct {
	s *state.State
	r *cmdroute.Router

	menus    []interfaces.Menu
	commands []discord.Command

	interfaces.App
}

func New(a interfaces.App) interfaces.CommandRepository {
	cr := &CommandRepositoryImpl{
		s: a.State(),
		r: a.CommandRouter(),

		App: a,
	}

	commands, err := cr.s.Commands(discord.AppID(cr.Config().Discord.AppIDSnowflake))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get commands from discord")
	}

	cr.commands = commands

	cr.menus = append(cr.menus, general.New(cr))
	cr.menus = append(cr.menus, bugs.New(cr))
	cr.menus = append(cr.menus, system.New(cr))

	return cr
}

func (cr *CommandRepositoryImpl) GetMenus() []interfaces.Menu {
	return cr.menus
}

func (cr *CommandRepositoryImpl) RegisterCommands() error {
	var commands []api.CreateCommandData

	for _, menu := range cr.menus {
		data, err := menu.RegisterCommands()
		if err != nil {
			return err
		}
		commands = append(commands, data...)
	}

	for _, command := range commands {
		logrus.WithField("command", command.Name).Info("Registering command")
	}

	_, err := cr.s.BulkOverwriteCommands(discord.AppID(cr.Config().Discord.AppIDSnowflake), commands)

	return err
}
