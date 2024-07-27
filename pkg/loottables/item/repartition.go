package item

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

func RepartitionUniform[I constraints.Integer](min, max I) func(*rand.Rand) I {
	return func(r *rand.Rand) I {
		return min + I(r.Intn(int(max-min+1)))
	}
}

func RepartitionGaussian[I constraints.Integer](mean, stdDev I) func(*rand.Rand) I {
	return func(r *rand.Rand) I {
		return mean + I(r.NormFloat64()*float64(stdDev))
	}
}

func RepartitionLinear[I constraints.Integer](min, max I) func(*rand.Rand) I {
	return func(r *rand.Rand) I {
		return min + I(r.Float64()*float64(max-min+1))
	}
}

func RepartitionExponential[I constraints.Integer](lambda float64) func(*rand.Rand) I {
	return func(r *rand.Rand) I {
		return I(r.ExpFloat64() / lambda)
	}
}
