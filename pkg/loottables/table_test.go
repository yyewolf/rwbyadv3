package loottables

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/yyewolf/rwbyadv3/internal/loots"
	"github.com/yyewolf/rwbyadv3/pkg/loottables/item"
)

func TestSimpleTable(t *testing.T) {
	simpleTable := New(
		item.New(&loots.Liens{}, 2).WithAmountRange(30, 150, 1),
	)

	r := rand.New(rand.NewSource(1))
	out := simpleTable.ChooseRandomItems(r, 2)

	// Verify that two items were returned
	if len(out) != 2 {
		t.Errorf("Expected 2 items, got %d", len(out))
	}
}

func TestAdvancedTable(t *testing.T) {
	advancedTable := New(
		item.New(&loots.Liens{}, 1).
			WithAmountRange(60, 150, 3).
			AlwaysDrop().
			OnlyDropOnce(),

		NewSubLootTable(1,
			AlwaysDrop(),
			Drop(2),
			WithEntries(
				item.New(&loots.Exit{}, 1).
					AlwaysDrop(),
				item.New(&loots.Liens{}, 1).
					WithAmountRange(160, 250, 1).
					WithRepartitionFunc(item.RepartitionGaussian(190, 10)),
				item.New(item.Nothing{}, 8),
			),
		),

		NewSubLootTable(1,
			AlwaysDrop(),
			Drop(1),
			WithEntries(
				item.New(&loots.Liens{}, 1).
					WithAmountRange(20, 100, 1).
					WithRepartitionFunc(item.RepartitionGaussian(50, 2)),
				item.New(item.Nothing{}, 64),
			),
		),
	)

	r := rand.New(rand.NewSource(5476584))
	var results [][]interface{}

	var testCount = 1000000

	for i := 0; i < testCount; i++ {
		out := advancedTable.ChooseRandomItems(r, 2)
		results = append(results, out)
	}

	var amountOfLiensPicked int
	for _, result := range results {
		// Check if Nothing was picked or Liens was picked
		liensPicked := false
		for _, item := range result {
			if item, ok := item.(*loots.Liens); ok {
				if item.Amount > 150 {
					liensPicked = true
				}
			}
		}
		if liensPicked {
			amountOfLiensPicked++
		}
	}

	fmt.Printf("Percentage of Liens picked: %0.2f%% \n", float64(amountOfLiensPicked*100)/float64(testCount))
}
