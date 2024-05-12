package rwbyadv3

import "embed"

//go:embed cards/img/*
var cardFS embed.FS

func GetCardFS() embed.FS {
	return cardFS
}
