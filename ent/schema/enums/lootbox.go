package enums

type LootBoxType string

const (
	Classic LootBoxType = "classic"
	Rare    LootBoxType = "rare"
	Special LootBoxType = "special"
	Limited LootBoxType = "limited"
)

// Values provides list valid values for Enum.
func (LootBoxType) Values() (kinds []string) {
	for _, s := range []LootBoxType{Classic, Rare, Special, Limited} {
		kinds = append(kinds, string(s))
	}
	return
}
