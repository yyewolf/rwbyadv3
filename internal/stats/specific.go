package stats

func HasChance(percent float64) bool {
	rolled := r.Float64()
	return rolled > percent
}
