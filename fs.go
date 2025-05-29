package rwbyadv3

import (
	"embed"
)

//go:embed cards/img/*
var cardFS embed.FS

func GetCardFS() embed.FS {
	return cardFS
}

//go:embed dungeons/dist/*
var dungeonFS embed.FS

func GetDungeonFS() embed.FS {
	return dungeonFS
}

//go:embed www/build/*
var wwwFS embed.FS

func GetWwwFS() embed.FS {
	return wwwFS
}
