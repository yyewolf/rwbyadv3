package interfaces

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Menu interface {
	GetName() string
	GetEmoji() discord.Emoji

	GetSubcommands() []Command
	RegisterCommands() ([]api.CreateCommandData, error)
}
