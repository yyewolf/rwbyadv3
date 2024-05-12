package cards

import (
	"fmt"
	"testing"
)

func TestParsing(t *testing.T) {
	ParseCards("../../cards/yml")

	for _, card := range Cards {
		fmt.Println(card)
	}
}
