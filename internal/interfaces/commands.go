package interfaces

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
)

type Command interface {
	GetName() string
	GetDescription() string
	RegisterCommand() (*api.CreateCommandData, error)

	Func() cmdroute.CommandHandlerFunc
}

type CommandRepository interface {
	RegisterCommands() error

	GetMenus() []Menu

	App
}
