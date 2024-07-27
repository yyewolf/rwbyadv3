package loots

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/pkg/loots/item"
)

type Liens struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Amount int    `json:"amount"`
}

func (m Liens) GetID() string {
	return m.ID
}

func (m Liens) GetType() string {
	return m.Type
}

func (m Liens) GetX() int {
	return m.X
}

func (m Liens) GetY() int {
	return m.Y
}

func (m Liens) Generate(r *rand.Rand, point [2]int) Loot {
	var buffer = make([]byte, 16)
	r.Read(buffer)
	m.ID = uuid.NewSHA1(uuid.NameSpaceDNS, buffer).String()

	m.Type = "money"
	m.X = point[0]
	m.Y = point[1]
	m.Amount = r.Intn(125) + 50
	return m
}

func (m Liens) PickedUp(tx *sql.Tx, p *models.Player) {
	p.Liens += int64(m.Amount)
}

func (m Liens) RewardText(l []Loot) string {
	amount := 0
	for _, loot := range l {
		if loot.GetType() == "money" {
			amount += loot.(Liens).Amount
		}
	}
	return fmt.Sprintf("You found **%d Ⱡ** (Liens)!", amount)
}

func (m *Liens) GetAmount() int {
	return m.Amount
}

func (m *Liens) SetAmount(amount int) {
	m.Amount = amount
}

func (m *Liens) New() item.Amountable[int] {
	return &Liens{}
}
