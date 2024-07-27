package loottables

import (
	"math/rand"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3/pkg/loottables/item"
)

// LootTable represents a list of LootTableEntries.
type LootTable []LootTableEntry

// LootTableEntry is an interface that all entries in the loot table should implement.
type LootTableEntry interface {
	GetWeight() float64
	SetWeight(weight float64)
	GetItem(r *rand.Rand) interface{}
	GetUnique() bool
	GetAlways() bool
	GetEnabled() bool
	SetEnabled(enabled bool)

	Copy() interface{}
}

func FormatResult(result []interface{}) []interface{} {
	var flatResult []interface{}

	for _, item := range result {
		if i, ok := item.([]interface{}); ok {
			flatResult = append(flatResult, i...)
		} else {
			flatResult = append(flatResult, item)
		}
	}

	for i := 0; i < len(flatResult); i++ {
		r := flatResult[i]
		// Remove item.Nothin
		if _, ok := r.(item.Nothing); ok {
			flatResult = append(flatResult[:i], flatResult[i+1:]...)
		}
	}

	return flatResult
}

func (lt LootTable) Copy() LootTable {
	var out LootTable
	for _, entry := range lt {
		if entry, ok := entry.Copy().(LootTableEntry); ok {
			out = append(out, entry)
		} else {
			logrus.WithField("entry", entry).Error("Failed to copy loot table entry")
		}
	}
	return out
}

// ChooseRandomItems selects a specified number of random items from the loot table.
func (lt LootTable) ChooseRandomItems(r *rand.Rand, amount int) []interface{} {
	// Deep copy the loot table
	lt = lt.Copy()

	var result []interface{}
	var alwaysDropItems []LootTableEntry

	// Collect always drop items and disable them for future consideration
	for _, entry := range lt {
		if entry.GetAlways() {
			alwaysDropItems = append(alwaysDropItems, entry)
		}
	}

	// Add always drop items to the result
	for _, item := range alwaysDropItems {
		if item.GetUnique() {
			// Drop it once
			result = append(result, item.GetItem(r))
			item.SetEnabled(false)
			if len(result) >= amount {
				return FormatResult(result)
			}
		} else {
			// Drop it multiple times
			for i := 0; i < int(item.GetWeight()); i++ {
				result = append(result, item.GetItem(r))
				if len(result) >= amount {
					return FormatResult(result)
				}
				item.SetWeight(item.GetWeight() - 1)
			}
		}
	}

	// Randomly select additional items if needed
	for len(result) < amount {
		totalWeight := 0.0
		for _, entry := range lt {
			if entry.GetEnabled() {
				totalWeight += entry.GetWeight()
			}
		}

		// If no items are left that can be selected, break the loop
		if totalWeight == 0 {
			break
		}

		randomPoint := r.Float64() * totalWeight
		for _, entry := range lt {
			if entry.GetEnabled() {
				randomPoint -= entry.GetWeight()
				if randomPoint <= 0 {
					item := entry.GetItem(r)
					result = append(result, item)
					if entry.GetUnique() {
						entry.SetEnabled(false)
					}
					entry.SetWeight(entry.GetWeight() - 1)
					break
				}
			}
		}
	}

	return FormatResult(result)
}

func New(entries ...LootTableEntry) LootTable {
	return entries
}
