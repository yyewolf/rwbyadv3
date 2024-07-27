package loottables

import (
	"math/rand"
)

type SubLootTable struct {
	Weight float64
	Count  int

	Unique  bool
	Always  bool
	Enabled bool

	LootTable
}

func (lt SubLootTable) GetWeight() float64 {
	totalWeight := 0.0
	for _, entry := range lt.LootTable {
		totalWeight += entry.GetWeight()
	}
	return totalWeight
}

func (lt *SubLootTable) SetWeight(weight float64) {
	lt.Weight = weight
}

func (lt SubLootTable) GetItem(r *rand.Rand) interface{} {
	return lt.LootTable.ChooseRandomItems(r, lt.Count)
}

func (lt SubLootTable) GetUnique() bool {
	return lt.Unique
}

func (lt SubLootTable) GetAlways() bool {
	return lt.Always
}

func (lt SubLootTable) GetEnabled() bool {
	return lt.Enabled
}

func (lt *SubLootTable) SetEnabled(enabled bool) {
	lt.Enabled = enabled
}

func (lt *SubLootTable) Copy() interface{} {
	out := SubLootTable{
		Weight:  lt.Weight,
		Count:   lt.Count,
		Unique:  lt.Unique,
		Always:  lt.Always,
		Enabled: lt.Enabled,
	}

	for _, entry := range lt.LootTable {
		if entry, ok := entry.Copy().(LootTableEntry); ok {
			out.LootTable = append(out.LootTable, entry)
		}
	}

	return &out
}

func NewSubLootTable(weight float64, opts ...func(*SubLootTable)) *SubLootTable {
	out := SubLootTable{Weight: weight, Enabled: true}

	for _, opt := range opts {
		opt(&out)
	}

	return &out
}

func WithEntry(entry LootTableEntry) func(*SubLootTable) {
	return func(lt *SubLootTable) {
		lt.LootTable = append(lt.LootTable, entry)
	}
}

func WithEntries(entries ...LootTableEntry) func(*SubLootTable) {
	return func(lt *SubLootTable) {
		lt.LootTable = append(lt.LootTable, entries...)
	}
}

func Unique() func(*SubLootTable) {
	return func(lt *SubLootTable) {
		lt.Unique = true
	}
}

func WithCount(count int) func(*SubLootTable) {
	return func(lt *SubLootTable) {
		lt.Count = count
	}
}

func Always() func(*SubLootTable) {
	return func(lt *SubLootTable) {
		lt.Always = true
	}
}
