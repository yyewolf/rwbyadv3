package stats

import (
	"math"
	"math/rand"
	"time"
)

type Dice interface {
	roll() float64
}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Roll[K Dice](r K) (result float64) {
	result = r.roll()
	for result < 0 {
		result += math.Abs(reconciliation.roll())
	}
	for result > 100 {
		result -= math.Abs(reconciliation.roll())
	}
	return
}
