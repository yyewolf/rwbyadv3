package database

import (
	"github.com/FrancoLiberali/cql"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yyewolf/rwbyadv3/internal/conditions"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
	"gorm.io/gorm"
)

type PlayerRepositoryImpl struct {
	*gorm.DB
}

func NewPlayerRepository(g *gorm.DB) interfaces.PlayerRepository {
	return &PlayerRepositoryImpl{
		DB: g,
	}
}

func (p *PlayerRepositoryImpl) Create(player *models.Player) error {
	result := p.DB.Create(player)

	return result.Error
}

func (p *PlayerRepositoryImpl) GetByDiscordID(discordID snowflake.ID) (*models.Player, error) {
	return cql.Query(
		p.DB,
		conditions.Player.DiscordID.Is().Eq(discordID.String()),
	).FindOne()
}
