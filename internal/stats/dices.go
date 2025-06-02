package stats

var (
	// Used when roll are under 0 or over 100 to keep them near their first value (-1 will likely be around 0)
	reconciliation Dice = Normal(0, 1)

	NormalLootDie  Dice = Normal(50, 15)
	RareLootDie    Dice = Normal(55, 15)
	LimitedLootDie Dice = Normal(62, 15)

	NormalLootExp  Dice = Exponential(0.1)
	RareLootExp    Dice = Exponential(0.09)
	LimitedLootExp Dice = Exponential(0.08)
)
