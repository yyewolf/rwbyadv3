package models

import (
	"github.com/yyewolf/rwbyadv3/internal/mongo"
)

type Database struct {
	*mongo.Database
}

func Migrate(d *mongo.Database) *Database {
	var db = &Database{
		Database: d,
	}

	return db
}
