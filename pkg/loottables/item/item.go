package item

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

type Amountable[I constraints.Integer] interface {
	SetAmount(I)
	GetAmount() I

	New(r *rand.Rand) Amountable[I]
}

// Item represents a specific item with a weight for the drop chance. And an amount, it can be used to represent money for example.
type Item[I constraints.Integer, T Amountable[I]] struct {
	Item T

	Weight float64

	// Used for cool looting mechanics
	Unique  bool
	Always  bool
	Enabled bool

	// Range of possible amounts
	Min I
	Max I

	// Either Step or RepartitionFunc should be set
	Step            I
	RepartitionFunc func(*rand.Rand) I
}

// GetWeight returns the weight of the item.
func (ie Item[I, T]) GetWeight() float64 {
	return ie.Weight
}

// SetWeight sets the weight of the item.
func (ie *Item[I, T]) SetWeight(weight float64) {
	ie.Weight = weight
}

// GetItem returns the actual item.
func (ie Item[I, T]) GetItem(r *rand.Rand) interface{} {
	item := ie.Item.New(r)

	if ie.Min == 0 && ie.Max == 0 {
		return item
	}

	if ie.RepartitionFunc != nil {
		item.SetAmount(ie.RepartitionFunc(r))
	} else {
		// Equal odds for each amount
		amounts := (ie.Max-ie.Min)/ie.Step + 1
		amount := ie.Min + I(r.Intn(int(amounts)))*ie.Step
		item.SetAmount(amount)
	}

	if item.GetAmount() < ie.Min {
		item.SetAmount(ie.Min)
	} else if item.GetAmount() > ie.Max {
		item.SetAmount(ie.Max)
	}

	return item
}

// GetUnique returns whether the item only drops once.
func (ie Item[I, T]) GetUnique() bool {
	return ie.Unique
}

// GetAlways returns whether the item always drops.
func (ie Item[I, T]) GetAlways() bool {
	return ie.Always
}

// GetEnabled returns whether the item is enabled.
func (ie Item[I, T]) GetEnabled() bool {
	return ie.Enabled
}

// SetEnabled sets whether the item is enabled.
func (ie *Item[I, T]) SetEnabled(enabled bool) {
	ie.Enabled = enabled
}

// Copy returns a copy of the item.
func (ie Item[I, T]) Copy() interface{} {
	return &Item[I, T]{
		Item:            ie.Item,
		Weight:          ie.Weight,
		Unique:          ie.Unique,
		Always:          ie.Always,
		Enabled:         ie.Enabled,
		Min:             ie.Min,
		Max:             ie.Max,
		Step:            ie.Step,
		RepartitionFunc: ie.RepartitionFunc,
	}
}

// New creates a new ItemWithAmount with the specified item and weight and step.
func New[I constraints.Integer, T Amountable[I]](item T, weight float64) *Item[I, T] {
	out := Item[I, T]{Item: item, Weight: weight, Enabled: true}

	return &out
}

func (i *Item[I, T]) WithRepartitionFunc(repartitionFunc func(*rand.Rand) I) *Item[I, T] {
	i.RepartitionFunc = repartitionFunc
	return i
}

func (i *Item[I, T]) WithAmountRange(min, max, step I) *Item[I, T] {
	i.Min = min
	i.Max = max
	i.Step = step
	return i
}

func (i *Item[I, T]) AlwaysDrop() *Item[I, T] {
	i.Always = true
	return i
}

func (i *Item[I, T]) OnlyDropOnce() *Item[I, T] {
	i.Unique = true
	return i
}
