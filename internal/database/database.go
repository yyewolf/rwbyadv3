package database

import (
	"fmt"

	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/internal/env"
	"github.com/yyewolf/rwbyadv3/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	c env.Config

	Players PlayerRepository

	db *gorm.DB
}

func New(c env.Config) *Database {
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

	return &Database{
		db: db,
		c:  c,

		Players: NewPlayerRepository(db),
	}
}

func (db *Database) Disconnect() error {
	sqlDB, err := db.db.DB()
	if err != nil {
		logrus.Fatal(err)
	}
	return sqlDB.Close()
}

func (db *Database) Migrate() error {
	return db.db.AutoMigrate(&models.Player{})
}
