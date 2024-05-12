package stats

func RandN(min, max int) int {
	return min + r.Intn(max-min)
}
