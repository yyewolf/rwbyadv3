package stats

var (
	// Used when roll are under 0 or over 100 to keep them near their first value (-1 will likely be around 0)
	reconciliation Dice = Normal(0, 1)

	NormalLootDie  Dice = Normal(48.5, 15)
	RareLootDie    Dice = Normal(55, 15)
	LimitedLootDie Dice = Normal(60, 15)

	NormalLootExp  Dice = Exponential(0.06)
	RareLootExp    Dice = Exponential(0.05)
	LimitedLootExp Dice = Exponential(0.046)
)
