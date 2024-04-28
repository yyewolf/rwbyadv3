package database

import (
	"github.com/FrancoLiberali/cql"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/yyewolf/rwbyadv3/internal/conditions"
	"github.com/yyewolf/rwbyadv3/models"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(player *models.Player) error
	GetByDiscordID(discordID discord.UserID) (*models.Player, error)
}

type PlayerRepositoryImpl struct {
	*gorm.DB
}

func NewPlayerRepository(g *gorm.DB) PlayerRepository {
	return &PlayerRepositoryImpl{
		DB: g,
	}
}

func (p *PlayerRepositoryImpl) Create(player *models.Player) error {
	result := p.DB.Create(player)

	return result.Error
}

func (p *PlayerRepositoryImpl) GetByDiscordID(discordID discord.UserID) (*models.Player, error) {
	return cql.Query(
		p.DB,
		conditions.Player.DiscordID.Is().Eq(discordID.String()),
	).FindOne()
}
