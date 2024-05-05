package stats

import (
	"fmt"
	"math"
	"testing"
)

var N = 100000

func rollForDice(d Dice) []float64 {
	// This draws the random number from the normal distribution to test the function
	data := make([]float64, N)
	for i := range data {
		data[i] = Roll(d)
	}
	return data
}

func showStatsNormal(data []float64) {
	fmt.Println("Under 10:", float64(countUnder(data, 10))*100.0/float64(N), "%")
	fmt.Println("Over 90:", float64(countOver(data, 90))*100.0/float64(N), "%")
	fmt.Println("Between 40-60:", float64(N-countOver(data, 60)-countUnder(data, 40))*100.0/float64(N), "%")
}

func showStatsExp(data []float64) {
	fmt.Println("Between 0-20:", float64(N-countOver(data, 20))*100.0/float64(N), "%")
	fmt.Println("Between 20-40:", float64(countOver(data, 20)-countOver(data, 40))*100.0/float64(N), "%")
	fmt.Println("Between 40-60:", float64(countOver(data, 40)-countOver(data, 60))*100.0/float64(N), "%")
	fmt.Println("Between 60-80:", float64(countOver(data, 60)-countOver(data, 80))*100.0/float64(N), "%")
	fmt.Println("Between 80-100:", float64(countOver(data, 80))*100.0/float64(N), "%")
}

func TestNormalLoot(t *testing.T) {
	data := rollForDice(NormalLootDie)
	showStatsNormal(data)
}

func TestRareLoot(t *testing.T) {
	data := rollForDice(RareLootDie)
	showStatsNormal(data)
}

func TestLimitedLoot(t *testing.T) {
	data := rollForDice(LimitedLootDie)
	showStatsNormal(data)
}

func TestNormalExp(t *testing.T) {
	data := rollForDice(NormalLootExp)
	showStatsExp(data)
}

func TestRareExp(t *testing.T) {
	data := rollForDice(RareLootExp)
	showStatsExp(data)
}

func TestLimitedExp(t *testing.T) {
	data := rollForDice(LimitedLootExp)
	showStatsExp(data)
}

func TestScenarioExp(t *testing.T) {
	normalN := int(math.Ceil(0.85 * float64(N)))
	rareN := int(math.Ceil(0.14 * float64(N)))
	limitedN := int(math.Ceil(0.01 * float64(N)))

	if normalN+rareN+limitedN != N {
		normalN -= normalN + rareN + limitedN - N
	}

	normalData := make([]float64, normalN)
	rareData := make([]float64, rareN)
	limitedData := make([]float64, limitedN)
	for i := range normalData {
		normalData[i] = Roll(NormalLootExp)
	}
	for i := range rareData {
		rareData[i] = Roll(RareLootExp)
	}
	for i := range limitedData {
		limitedData[i] = Roll(LimitedLootExp)
	}

	showStatsExp(append(append(normalData[:int(normalN)], rareData[:int(rareN)]...), limitedData[:int(limitedN)]...))
}
