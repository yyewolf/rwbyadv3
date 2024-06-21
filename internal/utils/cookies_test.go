package utils

import (
	"fmt"
	"testing"
)

func TestCookieToken(t *testing.T) {
	fmt.Println(GenerateNewCookieId())
	fmt.Println(GenerateNewCookieId())
	fmt.Println(GenerateNewCookieId())
	fmt.Println(GenerateNewCookieId())
	fmt.Println(GenerateNewCookieId())
}
