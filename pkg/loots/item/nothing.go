package item

type Nothing struct{}

func (m Nothing) GetAmount() int {
	return 0
}

func (m Nothing) SetAmount(amount int) {}

func (m Nothing) New() Amountable[int] {
	return Nothing{}
}
