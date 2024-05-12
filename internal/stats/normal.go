package stats

type NormalParams struct {
	Mean   float64
	StdDev float64
}

func Normal(mean, stdev float64) NormalParams {
	return NormalParams{
		Mean:   mean,
		StdDev: stdev,
	}
}

func (n NormalParams) roll() (result float64) {
	return r.NormFloat64()*n.StdDev + n.Mean
}
