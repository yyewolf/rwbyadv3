package loots

import (
	"database/sql"
	"math/rand"

	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/pkg/loots/item"
)

type Exit struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

func (m Exit) GetID() string {
	return m.ID
}

func (m Exit) GetType() string {
	return m.Type
}

func (m Exit) GetX() int {
	return m.X
}

func (m Exit) GetY() int {
	return m.Y
}

func (m Exit) Generate(r *rand.Rand, point [2]int) Loot {
	var buffer = make([]byte, 16)
	r.Read(buffer)
	m.ID = uuid.NewSHA1(uuid.NameSpaceDNS, buffer).String()

	m.Type = "exit"
	m.X = point[0]
	m.Y = point[1]
	return m
}

func (m Exit) PickedUp(tx *sql.Tx, p *models.Player) {
	// Do nothing
}

func (m Exit) RewardText(l []Loot) string {
	return ""
}

func (m *Exit) GetAmount() int {
	return 1
}

func (m *Exit) SetAmount(amount int) {}

func (m *Exit) New() item.Amountable[int] {
	return &Exit{}
}
