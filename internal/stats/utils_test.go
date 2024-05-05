package stats

func countUnder(data []float64, threshold float64) int {
	count := 0
	for _, v := range data {
		if v < threshold {
			count++
		}
	}
	return count
}

func countOver(data []float64, threshold float64) int {
	count := 0
	for _, v := range data {
		if v > threshold {
			count++
		}
	}
	return count
}
