package loots

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/yyewolf/rwbyadv3/models"
	"github.com/yyewolf/rwbyadv3/pkg/loottables/item"
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

func (m Liens) Place(point [2]int) DungeonLoot {
	m.X = point[0]
	m.Y = point[1]
	return m
}

func (m Liens) PickedUp(tx *sql.Tx, p *models.Player) {
	p.Liens += int64(m.Amount)
}

func (m Liens) RewardText(l []interface{}) string {
	amount := 0
	for _, loot := range l {
		switch loot := loot.(type) {
		case *Liens:
			amount += loot.Amount
		case Liens:
			amount += loot.Amount
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

func (m *Liens) New(r *rand.Rand) item.Amountable[int] {
	var buffer = make([]byte, 16)
	r.Read(buffer)
	m.ID = uuid.NewSHA1(uuid.NameSpaceDNS, buffer).String()

	m.Type = "money"
	return &Liens{
		ID:   m.ID,
		Type: m.Type,
	}
}
