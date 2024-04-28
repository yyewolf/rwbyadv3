package models

import (
	"github.com/FrancoLiberali/cql/model"
)

type Player struct {
	model.UUIDModel

	DiscordID string `gorm:"index,unique"`
}
