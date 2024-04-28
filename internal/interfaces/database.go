package interfaces

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/yyewolf/rwbyadv3/models"
)

type Database interface {
	Players() PlayerRepository

	Disconnect() error
	Migrate() error
}

type PlayerRepository interface {
	Create(player *models.Player) error
	GetByDiscordID(discordID discord.UserID) (*models.Player, error)
}
