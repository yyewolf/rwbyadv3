package utils

import "strings"

func Joinln(s ...string) string {
	return strings.Join(s, "\n")
}

func Optional[K any](e K) *K {
	return &e
}
