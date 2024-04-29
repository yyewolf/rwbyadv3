package interfaces

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/models"
)

type Database interface {
	Players() PlayerRepository

	Disconnect() error
	Migrate() error
}

type PlayerRepository interface {
	Create(player *models.Player) error
	GetByDiscordID(discordID snowflake.ID) (*models.Player, error)
}
