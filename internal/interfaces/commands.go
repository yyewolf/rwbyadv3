package interfaces

import (
	"github.com/diamondburned/arikawa/v3/api"
)

type Command interface {
	GetName() string
	GetDescription() string
	RegisterCommand() (*api.CreateCommandData, error)

	// This is not useful and can lead to issue with subcommands
	// Func() cmdroute.CommandHandlerFunc
}

type CommandRepository interface {
	RegisterCommands() error

	GetMenus() []Menu

	App
}
