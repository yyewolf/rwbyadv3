package utils

import "github.com/disgoorg/disgo/discord"

// HELPER FUNCTIONS
func FindCommandByName(commands []discord.ApplicationCommand, name string) discord.ApplicationCommand {
	for _, command := range commands {
		if command.Name() == name {
			return command
		}
	}

	return nil
}
