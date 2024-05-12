package stats

type ExponentialParams struct {
	Rate float64
}

func Exponential(rate float64) ExponentialParams {
	return ExponentialParams{
		Rate: rate,
	}
}

func (e ExponentialParams) roll() (result float64) {
	return r.ExpFloat64() / e.Rate
}
