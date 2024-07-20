package loots

import (
	"database/sql"
	"math/rand"

	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/models"
)

type MoneyBag struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Amount int    `json:"amount"`
}

func (m MoneyBag) GetID() string {
	return m.ID
}

func (m MoneyBag) GetType() string {
	return m.Type
}

func (m MoneyBag) GetX() int {
	return m.X
}

func (m MoneyBag) GetY() int {
	return m.Y
}

func (m MoneyBag) Generate(r *rand.Rand, point [2]int) Loot {
	var buffer = make([]byte, 16)
	r.Read(buffer)
	m.ID = uuid.NewSHA1(uuid.NameSpaceDNS, buffer).String()

	m.Type = "money"
	m.X = point[0]
	m.Y = point[1]
	m.Amount = r.Intn(125) + 50
	return m
}

func (m MoneyBag) PickedUp(tx *sql.Tx, p *models.Player) {
	p.Liens += int64(m.Amount)
}
