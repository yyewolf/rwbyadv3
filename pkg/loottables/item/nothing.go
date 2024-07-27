package item

import "math/rand"

type Nothing struct{}

func (m Nothing) GetAmount() int {
	return 0
}

func (m Nothing) SetAmount(amount int) {}

func (m Nothing) New(r *rand.Rand) Amountable[int] {
	return Nothing{}
}
