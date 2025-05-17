package enums

type LootBoxType string

const (
	LootBoxClassic LootBoxType = "classic"
	LootBoxRare    LootBoxType = "rare"
	LootBoxSpecial LootBoxType = "special"
	LootBoxLimited LootBoxType = "limited"
)

func LootBoxTypes() []LootBoxType {
	return []LootBoxType{LootBoxClassic, LootBoxRare, LootBoxSpecial, LootBoxLimited}
}

// Values provides list valid values for Enum.
func (LootBoxType) Values() (kinds []string) {
	for _, s := range []LootBoxType{LootBoxClassic, LootBoxRare, LootBoxSpecial, LootBoxLimited} {
		kinds = append(kinds, string(s))
	}
	return
}
