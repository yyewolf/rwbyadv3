package database

import (
	"fmt"

	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/internal/interfaces"
	"github.com/yyewolf/rwbyadv3/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseImpl struct {
	c env.Config

	players interfaces.PlayerRepository

	db *gorm.DB
}

func New(c env.Config) interfaces.Database {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.User,
		c.Database.Pass,
		c.Database.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		panic("failed to connect database")
	}

	return &DatabaseImpl{
		db: db,
		c:  c,

		players: NewPlayerRepository(db),
	}
}

func (db *DatabaseImpl) Disconnect() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		logrus.Fatal(err)
	}
	return sqlDB.Close()
}

func (db *DatabaseImpl) Migrate() error {
	return db.db.AutoMigrate(&models.Player{})
}

func (db *DatabaseImpl) Players() interfaces.PlayerRepository {
	return db.players
}
