package utils

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

// HELPER FUNCTIONS
func FindCommandByName(commands []discord.Command, name string) discord.Command {
	for _, command := range commands {
		if command.Name == name {
			return command
		}
	}

	return discord.Command{}
}
